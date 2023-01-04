package arguments

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	extools "src/exeiac/tools"

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
		NonInteractive   bool     `yaml:"non_interactive"`
		BricksSpecifiers []string `yaml:"bricks_specifiers"`
		OtherOptions     []string `yaml:"other_options"`
	} `yaml:"default_arguments"`
}

// A structure matching the `Arguments` type one. It is used to merge the values defined both
// in the command lien with flags, and in the exeiac configuration file.
//
// NOTE(half-shell): Should be able to embed the ` Arguments` struct.
// We replicate the `Arguments` since it does not seem to work out of the box for some reason.
type Configuration struct {
	Action            string
	BricksNames       []string
	BricksSpecifiers  []string
	Interactive       bool
	Format            string
	Modules           map[string]string
	OtherOptions      []string
	Rooms             map[string]string
	ConfigurationFile string
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

func CreateConfiguration(confFilePath string) (configuration Configuration, err error) {
	file, err := os.ReadFile(confFilePath)
	if err != nil {
		return
	}

	var confFile ConfigurationFile
	err = yaml.Unmarshal(file, &confFile)
	if err != nil {
		return
	}

	rooms := make(map[string]string)
	modules := make(map[string]string)

	for _, room := range confFile.Rooms {
		rooms[room.Name] = room.Path
	}

	for _, module := range confFile.Modules {
		modules[module.Name] = module.Path
	}

	configuration = Configuration{
		Rooms:            rooms,
		Modules:          modules,
		Interactive:      confFile.DefaultArgs.NonInteractive,
		BricksSpecifiers: confFile.DefaultArgs.BricksSpecifiers,
		OtherOptions:     confFile.DefaultArgs.OtherOptions,
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

	if args.ConfigurationFile != "" {
		conf, err = CreateConfiguration(args.ConfigurationFile)
	} else {
		var configFilePath string

		configFilePath, err = xdg.SearchConfigFile(CONFIG_FILE)
		if err == nil {
			conf, err = CreateConfiguration(configFilePath)
		}
	}

	if err != nil {
		// NOTE(half-shell): We only report an error on configuration reading if the command line
		// arguments are enough to handle exeiac's execution.
		if !(len(args.Action) > 0 && len(args.BricksNames) > 0 && len(args.Rooms) > 0) {
			return
		}
	}

	modules := make(map[string]string)
	rooms := make(map[string]string)

	if err == nil {
		modules = conf.Modules
		rooms = conf.Rooms
	} else {
		// NOTE(half-shell): We avoid propagating the error up the call stack
		// since we're handling it.
		err = nil
	}

	for name, path := range args.Modules {
		modules[name] = path
	}

	for name, path := range args.Rooms {
		rooms[name] = path
	}

	// NOTE(half-shell): Ideally, we wouldn't want to mix exeiac injected flags and the
	// ones provided by the user. However, distinction is not needed for now.
	var other_options = append(conf.OtherOptions, args.OtherOptions...)
	if args.NonInteractive {
		other_options = extools.Deduplicate(append(other_options, "--non-interactive"))
	}

	configuration = Configuration{
		Action:            args.Action,
		BricksNames:       args.BricksNames,
		BricksSpecifiers:  extools.Deduplicate(append(conf.BricksSpecifiers, args.BricksSpecifiers...)),
		ConfigurationFile: args.ConfigurationFile,
		Format:            args.Format,
		Interactive:       (conf.Interactive && !args.NonInteractive) || args.Interactive,
		Modules:           modules,
		Rooms:             rooms,
		OtherOptions:      other_options,
	}

	return
}
