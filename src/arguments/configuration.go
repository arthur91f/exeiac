package arguments

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	extools "src/exeiac/tools"

	"github.com/PaesslerAG/jsonpath"
	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

const CONFIG_FILE = "exeiac/exeiac.yml"

// A type matching the structure of exeiac's configuration file (`conf.yml`)
// Its purpose is merely to load the configuration.
type ConfigurationFile struct {
	Rooms []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
	} `yaml:"rooms,flow"`
	Modules []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
	} `yaml:"modules,flow"`
	DefaultArgs struct {
		NonInteractive       bool      `yaml:"non_interactive"`
		BricksSpecifiers     []string  `yaml:"bricks_specifiers"`
		OtherOptions         []string  `yaml:"other_options"`
		InputNeededFor       *[]string `yaml:"inputs_needed_for"`       // NOTE1: why use pointer: https://stackoverflow.com/questions/54765161/how-to-test-for-default-values-in-golang-when-yaml-has-already-initialized-the-v
		InputNotNeededFor    *[]string `yaml:"inputs_not_needed_for"`   // see NOTE1
		DefaultIsInputNeeded *bool     `yaml:"default_is_input_needed"` // see NOTE1
	} `yaml:"default_arguments"`
}

type Room struct {
	Path string
	Name string
}

// A structure matching the `Arguments` type one. It is used to merge the values defined both
// in the command lien with flags, and in the exeiac configuration file.
//
// NOTE(half-shell): Should be able to embed the ` Arguments` struct.
// We replicate the `Arguments` since it does not seem to work out of the box for some reason.
type Configuration struct {
	Action                 string
	BricksNames            []string
	JsonPath               string
	BricksSpecifiers       []string
	Format                 string
	Interactive            bool
	ExceptionIsInputNeeded []string
	DefaultIsInputNeeded   bool
	Modules                map[string]string
	OtherOptions           []string
	Rooms                  []Room
	ConfigurationFilePath  string
}

func (a Configuration) String() string {
	var sb strings.Builder

	v := reflect.ValueOf(a)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		sb.WriteString("\t")
		sb.WriteString(t.Field(i).Name)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%v", v.Field(i).Interface()))
		sb.WriteString("\n")
	}

	return sb.String()
}

func CreateConfiguration(configurationFilePath string) (configuration Configuration, err error) {
	file, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return
	}

	var confFile ConfigurationFile
	err = yaml.Unmarshal(file, &confFile)
	if err != nil {
		return
	}

	rooms := make([]Room, len(confFile.Rooms))
	modules := make(map[string]string)

	for index, room := range confFile.Rooms {
		rooms[index] = Room{
			Name: room.Name,
			Path: room.Path,
		}
	}

	for _, module := range confFile.Modules {
		modules[module.Name] = module.Path
	}

	// input need options
	var defaultIsInputNeeded bool
	var exceptionIsInputNeeded []string
	if confFile.DefaultArgs.DefaultIsInputNeeded != nil {
		defaultIsInputNeeded = *confFile.DefaultArgs.DefaultIsInputNeeded
	} else {
		defaultIsInputNeeded = false
	}

	if defaultIsInputNeeded {
		if confFile.DefaultArgs.InputNotNeededFor != nil {
			exceptionIsInputNeeded = *confFile.DefaultArgs.InputNotNeededFor
		} else {
			exceptionIsInputNeeded = []string{"output", "init", "validate_code", "help"} // show, clean and get-depends aren't module action
		}
	} else {
		if confFile.DefaultArgs.InputNeededFor != nil {
			exceptionIsInputNeeded = *confFile.DefaultArgs.InputNeededFor
		} else {
			exceptionIsInputNeeded = []string{"plan", "lay", "remove"}
		}
	}

	configuration = Configuration{
		ConfigurationFilePath:  configurationFilePath,
		Rooms:                  rooms,
		Modules:                modules,
		Interactive:            confFile.DefaultArgs.NonInteractive,
		ExceptionIsInputNeeded: exceptionIsInputNeeded,
		DefaultIsInputNeeded:   defaultIsInputNeeded,
		BricksSpecifiers:       confFile.DefaultArgs.BricksSpecifiers,
		OtherOptions:           confFile.DefaultArgs.OtherOptions,
	}

	return
}

// Creates a `Configuration` struct resulting from the merger of the configuration file, and
// the command lien arguments (`Arguments`).
// Takes an instance of the `Arguments` struct as input that it can use to load the configuration file.
// Returns a tuple of a configuration if successful, return an error as last member otherwise.
//
// NOTE(half-shell): We can change the behaviour of the configuration building
// depending on a flag defined in "Arguments".
// For instance: do we want to merge arguments or override them?
// Current behaviour is "merging" them
func FromArguments(args Arguments) (configuration Configuration, err error) {
	var conf Configuration
	var configurationFilePath string

	if args.ConfigurationFilePath != "" {
		if filepath.IsAbs(args.ConfigurationFilePath) {
			configurationFilePath = args.ConfigurationFilePath
		} else {
			configurationFilePath, err = filepath.Abs(args.ConfigurationFilePath)
		}
	} else {
		configurationFilePath, err = xdg.SearchConfigFile(CONFIG_FILE)
		if err != nil {
			return
		}
	}

	conf, err = CreateConfiguration(configurationFilePath)

	if err != nil {
		// NOTE(half-shell): We only report an error on configuration reading if the command line
		// arguments are enough to handle exeiac's execution.
		if !(len(args.Action) > 0 && len(args.BricksNames) > 0 && len(args.Rooms) > 0) {
			return
		}
	}

	modules := make(map[string]string)
	rooms := conf.Rooms

	if err == nil {
		for name, path := range conf.Modules {
			var absPath string

			if filepath.IsAbs(path) {
				absPath = path
			} else {
				if !args.ListBricks {
					fmt.Fprintf(os.Stderr,
						"Warning: module path \"%s\" for module \"%s\" is relative. "+
							"Favor using an absolute path.\n", path, name)
				}

				absPath = filepath.Join(filepath.Dir(conf.ConfigurationFilePath), path)
			}

			modules[name] = absPath
		}

		for _, room := range conf.Rooms {
			var absPath string

			if filepath.IsAbs(room.Path) {
				absPath = room.Path
			} else {
				if !args.ListBricks {
					fmt.Fprintf(os.Stderr,
						"Warning: room path \"%s\" for room \"%s\" is relative. "+
							"Favor using an absolute path.\n",
						room.Path, room.Name)
				}

				absPath = filepath.Join(filepath.Dir(conf.ConfigurationFilePath), room.Path)
			}

			rooms = append(rooms, Room{
				Path: absPath,
				Name: room.Name,
			})
		}
	} else {
		// NOTE(half-shell): We avoid propagating the error up the call stack
		// since we're handling it.
		err = nil
	}

	for name, path := range args.Modules {
		var absPath string

		if filepath.IsAbs(path) {
			absPath = path
		} else {
			absPath = filepath.Join(filepath.Dir(conf.ConfigurationFilePath), path)
		}

		modules[name] = absPath
	}

	for name, path := range args.Rooms {
		var absPath string

		if filepath.IsAbs(path) {
			absPath = path
		} else {
			absPath = filepath.Join(filepath.Dir(conf.ConfigurationFilePath), path)
		}

		rooms = append(rooms, Room{
			Path: absPath,
			Name: name,
		})
	}

	_, err = jsonpath.New(args.JsonPath)
	if err != nil {
		return configuration, fmt.Errorf("%s is not valid", args.JsonPath)
	}

	// NOTE(half-shell): Ideally, we wouldn't want to mix exeiac injected flags and the
	// ones provided by the user. However, distinction is not needed for now.
	var other_options = append(conf.OtherOptions, args.OtherOptions...)
	if args.NonInteractive {
		other_options = extools.Deduplicate(append(other_options, "--non-interactive"))
	}

	configuration = Configuration{
		Action:                 args.Action,
		BricksNames:            args.BricksNames,
		BricksSpecifiers:       extools.Deduplicate(append(conf.BricksSpecifiers, args.BricksSpecifiers...)),
		Format:                 args.Format,
		Interactive:            (conf.Interactive && !args.NonInteractive) || args.Interactive,
		Modules:                modules,
		Rooms:                  rooms,
		OtherOptions:           other_options,
		ConfigurationFilePath:  configurationFilePath,
		JsonPath:               args.JsonPath,
		ExceptionIsInputNeeded: conf.ExceptionIsInputNeeded,
		DefaultIsInputNeeded:   conf.DefaultIsInputNeeded,
	}

	return
}
