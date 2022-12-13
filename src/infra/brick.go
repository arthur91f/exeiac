package infra

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	extools "src/exeiac/tools"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"gopkg.in/yaml.v2"
)

type Input struct {
	// The name the variable is supposed to take
	VarName string
	// The JSON path to access the variable
	JsonPath string
	// A reference to the related brick
	Brick *Brick
	// The type can be env_var, env_file, json, yaml, hashicorp
	Type InputFormat
	// The relative path from the brickPath of the file where the input will be written
	Path string // (obviously it is "" for env_var type)
}

func (i Input) String() string {
	return fmt.Sprintf("%s(%s):%s -> %s:%v",
		i.Path, i.Type, i.VarName, i.Brick.Name, i.JsonPath)
}

type Brick struct {
	// The brick's index. It represents the absolute brick ordering
	Index int
	// The brick's name. Usually the name of the parent directory
	Name string
	// The absolute path of the brick's directory
	Path string
	// The absolute path of the `brick.yml` file
	ConfigurationFilePath string
	// Wheither or not the brick contains a `brick.yml` file.
	// Meaning it does not contain any other brick.
	IsElementary bool
	// A pointer to a module
	Module *Module
	// Pointer to the bricks it depends on
	DirectPrevious Bricks

	// Data (and their represenation) from other bricks output that brick need to plan,lay,remove,output
	Inputs []Input

	Output []byte
	// Error from the last call to `Enrich()`
	EnrichError error
}

type BrickConfYaml struct {
	// The configuration file's format version
	Version string `yaml:"version"`
	// The name or path of the module it uses
	Module string `yaml:"module"`
	// A slice of different kinds of input needed for this brick
	// It **usually** matches the plain or processed output of another brick
	Input []struct {
		// The type of input this brick is expecting
		// Can match the strings "env_file" or "env_var"
		Type string `yaml:"type"`
		// If the type is a path, it is the path the dependency output should be saved to
		Path string `yaml:"path"`
		Data []struct {
			// The name the variable is expected to have
			Name string `yaml:"name"`
			// The key path the input variable should match
			// It is of the form "<brick_name>:"<json_path>"
			// OR "<brick_path>:"<json_path>"
			// e.g. "super-brick/brick:.object.field
			From string `yaml:"from"`
		}
	} `yaml:"input"`
}

// Reads a the yaml configuration file and parses it.
// Returns the parsed configuration.
// Returns an error because of reading or parsing the file.
func (bcy BrickConfYaml) New(path string) (BrickConfYaml, error) {
	in, err := os.ReadFile(path)
	if err != nil {
		return BrickConfYaml{}, err
	}

	err = yaml.Unmarshal(in, &bcy)
	if err != nil {
		return BrickConfYaml{}, err
	}

	return bcy, nil
}

func (b Brick) String() string {
	alwaysPresent := fmt.Sprintf(
		"index: %d\nname: %s\npath: %s\nisElementary: %t\nconfFile: %s",
		b.Index, b.Name, b.Path, b.IsElementary, b.ConfigurationFilePath)

	conditional := "\n"
	if b.EnrichError != nil {
		conditional = fmt.Sprintf("enrichError:%s\n", b.EnrichError)
	}

	if b.Module != nil {
		conditional = fmt.Sprintf("%smodule:%s\n", conditional, b.Module.Name)
	}

	if len(b.Inputs) > 0 {
		dpStr := []string{}
		for _, d := range b.Inputs {
			dpStr = append(dpStr, d.String())
		}
		conditional = fmt.Sprintf("%sinputData:%s",
			conditional, extools.StringListOfString(dpStr))
	}

	return fmt.Sprintf("%s%s", alwaysPresent, conditional)
}

func (brick *Brick) SetElementary(cfp string) *Brick {
	brick.IsElementary = true
	brick.ConfigurationFilePath = cfp

	return brick
}

// Processes the relevant parts of a brick's configuration and updates the brick itself
// with it.
func (brick *Brick) Enrich(bcy BrickConfYaml, infra *Infra) error {
	if !brick.IsElementary {
		return errors.New("Cannot enrich a non-elementary brick")
	}

	module, err := GetModule(bcy.Module, &infra.Modules)
	if err != nil {
		return err
	}

	brick.Module = module
	dependencies, err := bcy.resolveDependencies(infra)
	if err != nil {
		log.Println("An error occured when getting dependencies: ", err)
	}

	brick.Inputs = dependencies

	for _, i := range brick.Inputs {
		brick.DirectPrevious = append(brick.DirectPrevious, i.Brick)
	}
	brick.DirectPrevious = RemoveDuplicates(brick.DirectPrevious)

	return nil
}

// Parses this brick's input brick dependencies JSON output, and creates a map of formatters.
// Returns a map with the intput file path as the key, and the relevant Formatter as the value.
// The key is `env` if there is no path and the inputs are supposed to be passed around as
// environment variables
func (b *Brick) CreateFormatters() (formatters map[string]Formatter, err error) {
	formatters = make(map[string]Formatter)

	// NOTE(half-shell): This is pretty ugly, but there does not seem to be a way to create
	// a 3-fold nested map otherwise.
	// This is a pretty good use case for a tree-like structure!
	varNameToVal := make(map[string]interface{})
	pathToVars := make(map[string]map[string]interface{})
	// Temporary variable holding the values dispatched by format, path and variable name.
	// e.g. rawInputs[<data_format>][<file_path>][<variable_name>] => <variable_value>
	rawInputs := make(map[InputFormat]map[string]map[string]interface{})

	for _, i := range b.Inputs {
		var output interface{}
		err := json.Unmarshal(i.Brick.Output, &output)
		if err != nil {
			log.Fatalf("Could not parse JSON: %v", err)
		}

		varVal, err := jsonpath.Get(i.JsonPath, output)
		if err != nil {
			log.Fatalf("Error happened when solving dependency JSON path %v: %v", i.JsonPath, err)
		}

		path := filepath.Join(b.Path, i.Path)
		varNameToVal[i.VarName] = varVal
		pathToVars[path] = varNameToVal
		rawInputs[i.Type] = pathToVars
	}

	for format, paths := range rawInputs {
		for path, vals := range paths {
			switch format {
			case Json:
				formatters[path] = JsonFormat(vals)
			case Env:
				formatters[path] = EnvFormat(vals)
			default:
				// TODO(half-shell): One way of dealing with inputs passed around as environment variables
				// would be to check for a path, and if none is present, just return on for a path of `env`
				// for instance, which would be handled some other way down the line
				return formatters, errors.New(fmt.Sprintf("Format %s is not handled", format))
			}
		}
	}

	return formatters, nil
}

func (bcy BrickConfYaml) resolveDependencies(infra *Infra) ([]Input, error) {
	var dependencies []Input
	parseFromField := func(from string) (brickName string, dataKey string) {
		fields := strings.Split(from, ":")

		return fields[0], fields[1]
	}

	for _, i := range bcy.Input {
		for _, d := range i.Data {
			brickName, keyPath := parseFromField(d.From)
			// NOTE(half-shell): We sanitize the brick name here in case they turn
			// out to be brick paths
			brick, ok := infra.Bricks[SanitizeBrickName(brickName)]

			if !ok {
				return dependencies, errors.New(fmt.Sprintf("No brick names %s", brickName))
			}

			// NOTE(half-shell): We make sure the jsonPath's form is valid
			_, err := jsonpath.New(keyPath)
			if err != nil {
				return dependencies, err
			}

			if inputType, isSupported := SupportedFormats[i.Type]; isSupported {
				dependencies = append(dependencies, Input{
					VarName:  d.Name,
					JsonPath: keyPath,
					Brick:    brick,
					Type:     inputType,
					Path:     i.Path,
				})
			} else {
				err = errors.New(fmt.Sprintf("Format %s not supported", i.Type))
			}
		}
	}

	return dependencies, nil
}
