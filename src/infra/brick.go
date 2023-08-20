package infra

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	extools "src/exeiac/tools"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"gopkg.in/yaml.v2"
)

type From struct {
	Brick    *Brick
	Source   string // can be output or event
	JsonPath string
}

type Dependency struct {
	// A reference to the related brick
	From               From
	TriggeredAction    []string
	TriggerType        string
	DefaultNeededFor   bool     // Is inputs needed for an action by default
	ExceptionNeededFor []string // For which action the input doesn't respect the default needed behaviour
	Value              any
}

type Input struct {
	VarName    string
	Dependency *Dependency
	Type       string      // can be env_var or file
	Format     InputFormat // can be env, json or yaml
	Path       string      // (obviously it is "" for env_var type)
}

func (i Input) String() string {
	var need_for_str string
	if i.Dependency.DefaultNeededFor {
		need_for_str = "not need for"
	} else {
		need_for_str = "need for"
	}
	return fmt.Sprintf("%s(%s):%s -> %s:%v\n\t\t  %s %s",
		i.Path, i.Type, i.VarName,
		i.Dependency.From.Brick.Name,
		i.Dependency.From.JsonPath,
		need_for_str,
		i.Dependency.ExceptionNeededFor)
}

func (d Dependency) IsDependencyNeeded(
	action string,
) bool {
	if extools.ContainsString(d.ExceptionNeededFor, action) {
		return !d.DefaultNeededFor
	}
	return d.DefaultNeededFor
}

func (i Input) IsInputNeeded(
	action string,
) bool {
	return i.Dependency.IsDependencyNeeded(action)
}

func (b *Brick) GetInputsThatCallthisOutput(
	brick *Brick,
	jsonpath string,
) (
	inputs []Input,
) {
	for _, i := range b.Inputs {
		if i.Dependency.From.Brick == brick {
			if extools.AreJsonPathsLinked(jsonpath, i.Dependency.From.JsonPath) {
				inputs = append(inputs, i)
			}
		}
	}
	return
}

func (i Input) StringCompact() string {
	return fmt.Sprintf("%s -> %s:%v",
		i.VarName, i.Dependency.From.Brick.Name, i.Dependency.From.JsonPath)
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

	// Dependencies, when they are needed and how to get their values
	Dependencies map[string]*Dependency

	// How to present dependencies value for executing brick
	Inputs []Input

	Output []byte
	// Error from the last call to `Enrich()`
	EnrichError error
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
func (b *Brick) SetElementary(cfp string) {
	b.IsElementary = true
	b.ConfigurationFilePath = cfp
}

// Parses this brick's input brick dependencies JSON output, and creates a map of formatters.
// Returns a map with the intput file path as the key, and the relevant Formatter as the value.
// The key is `env` if there is no path and the inputs are supposed to be passed around as
// environment variables
func (b *Brick) CreateFormatters(
	action string,
) (
	fileFormatters map[string]Formatter, env_formatters EnvFormat, err error,
) {
	fileFormatters = make(map[string]Formatter)

	// Temporary variable holding the values dispatched by format, path and variable name.
	// e.g. rawInputs[<data_format>][<file_path>][<variable_name>] => <variable_value>
	// NOTE: This is a pretty good use case for a tree-like structure!
	rawInputs := make(map[InputFormat]map[string]map[string]interface{})

	var neededInputs []Input
	for _, i := range b.Inputs {
		if i.IsInputNeeded(action) {
			neededInputs = append(neededInputs, i)
		}
	}

	for _, i := range neededInputs {
		var output interface{}
		err := json.Unmarshal(i.Dependency.From.Brick.Output, &output)
		if err != nil {
			log.Fatalf("Could not parse JSON that correspond to %s output: %v", i.Dependency.From.Brick.Name, err)
		}

		varVal, err := jsonpath.Get(i.Dependency.From.JsonPath, output)
		if err != nil {
			log.Fatalf("Error happened when solving dependency JSON path of brick %v (jsonPath: '%v'): %v",
				b.Name, i.Dependency.From.JsonPath, err)
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

func ParseOutputName(from string) (brickName string, dataKey string, err error) {
	if from == "" {
		return "", "", fmt.Errorf("EMPTY")
	}

	fields := strings.Split(from, ":")

	if len(fields) != 2 {
		return "", "", fmt.Errorf("BAD FORMAT")
	}

	return fields[0], fields[1], nil
}
