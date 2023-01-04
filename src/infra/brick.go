package infra

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	// The Format can be env, json, yaml, hashicorp
	Format InputFormat
	// The type can be env, file
	Type string
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
	// The brick pointer of the room. Is set to nil if the brick is a room
	Room *Brick
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
		// Can match the strings "file" or "env_vars"
		Type string `yaml:"type"`
		// Can be json, yaml, env
		Format string `yaml:"format"`
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
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\tIndex: %d\n", b.Index))
	sb.WriteString(fmt.Sprintf("\tName: %s\n", b.Name))
	sb.WriteString(fmt.Sprintf("\tPath: %s\n", b.Path))
	sb.WriteString(fmt.Sprintf("\tIsElementary: %v", b.IsElementary))

	if len(b.ConfigurationFilePath) != 0 {
		sb.WriteString(fmt.Sprintf("\n\tConfigurationFile: %s", b.ConfigurationFilePath))
	}

	if b.EnrichError != nil {
		sb.WriteString(fmt.Sprintf("\n\tEnrichError:%s", b.EnrichError))
	}

	if b.Module != nil {
		sb.WriteString(fmt.Sprintf("\n\tModule:%s", b.Module.Name))
	}

	if len(b.Inputs) > 0 {
		sb.WriteString(fmt.Sprintf("\n\tInputs:"))
		for _, input := range b.Inputs {
			sb.WriteString(fmt.Sprintf("\n\t\t- %v", input))
		}
	}

	sb.WriteString("\n")

	return sb.String()
}

// Set's a brick as an elementary brick by setting its `IsElementary` flag
// to `true` and its `ConfigurationFilePath` with the provided one.
func (brick *Brick) SetElementary(cfp string) {
	brick.IsElementary = true
	brick.ConfigurationFilePath = cfp
}

// Processes the relevant parts of a brick's configuration and updates the brick itself
// with it.
func (brick *Brick) Enrich(bcy BrickConfYaml, infra *Infra) error {
	if !brick.IsElementary {
		return errors.New("Cannot enrich a non-elementary brick")
	}

	module, err := infra.GetModule(bcy.Module, brick)
	if err != nil {
		return err
	}

	brick.Module = module
	dependencies, err := bcy.resolveDependencies(infra)
	if err != nil {
		log.Printf("An error occured when getting dependencies of brick %s: %v\n", brick.Name, err)
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
func (b *Brick) CreateFormatters() (fileFormatters map[string]Formatter, env_formatters EnvFormat, err error) {
	fileFormatters = make(map[string]Formatter)

	// Temporary variable holding the values dispatched by format, path and variable name.
	// e.g. rawInputs[<data_format>][<file_path>][<variable_name>] => <variable_value>
	// NOTE: This is a pretty good use case for a tree-like structure!
	rawInputs := make(map[InputFormat]map[string]map[string]interface{})

	for _, i := range b.Inputs {
		var output interface{}
		err := json.Unmarshal(i.Brick.Output, &output)
		if err != nil {
			log.Fatalf("Could not parse JSON that correspond to %s output: %v", i.Brick.Name, err)
		}

		varVal, err := jsonpath.Get(i.JsonPath, output)
		if err != nil {
			log.Fatalf("Error happened when solving dependency JSON path of brick %v (jsonPath: '%v'): %v", b.Name, i.JsonPath, err)
		}

		path := filepath.Join(b.Path, i.Path)
		if _, exist := rawInputs[i.Format]; !exist {
			rawInputs[i.Format] = make(map[string]map[string]interface{})
		}
		if _, exist := rawInputs[i.Format][path]; !exist {
			rawInputs[i.Format][path] = make(map[string]interface{})
		}
		rawInputs[i.Format][path][i.VarName] = varVal
	}

	for format, paths := range rawInputs {
		for path, vals := range paths {
			switch format {
			case Json:
				fileFormatters[path] = JsonFormat(vals)
			case Env:
				if path == b.Path {
					env_formatters = EnvFormat(vals)
				} else {
					fileFormatters[path] = EnvFormat(vals)
				}
			default:
				// TODO(half-shell): One way of dealing with inputs passed around as environment variables
				// would be to check for a path, and if none is present, just return on for a path of `env`
				// for instance, which would be handled some other way down the line
				err = fmt.Errorf("Format %s is not handled", format)

				return
			}
		}
	}
	return
}

// Loops through a `Brick`'s configuration file's `Input` and builds a slice of `Input`s
// out of it.
// The `infra` argument is used to resolve a brick's name to a `Brick` reference.
// Returns an error if something wrong happened during `Brick's` name-to-reference resolution,
// or when checking that the `JsonPath` is a valid JSONPath.
func (bcy BrickConfYaml) resolveDependencies(infra *Infra) (inputs []Input, err error) {
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
				err = fmt.Errorf("No brick names %s", brickName)

				return
			}

			// NOTE(half-shell): We make sure the jsonPath's form is valid
			_, err = jsonpath.New(keyPath)
			if err != nil {
				return
			}

			if inputFormat, isSupported := SupportedFormats[i.Format]; isSupported {
				inputs = append(inputs, Input{
					VarName:  d.Name,
					JsonPath: keyPath,
					Brick:    brick,
					Format:   inputFormat,
					Type:     i.Type,
					Path:     i.Path,
				})
			} else {
				err = errors.New(fmt.Sprintf("Format %s not supported in %s",
					i.Type,
					brick.ConfigurationFilePath))
			}
		}
	}

	return
}
