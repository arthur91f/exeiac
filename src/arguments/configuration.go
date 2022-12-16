package arguments

import (
	"fmt"
	"os"
	"reflect"
	extools "src/exeiac/tools"
	"strings"

	"gopkg.in/yaml.v2"
)

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

func createConfiguration(confFilePath string) (configuration Configuration, err error) {
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

// NOTE(half-shell): We can change the behaviour of the configuration building
// depending on a flag defined in "Arguments".
// For instance: do we want to merge arguments or override them?
// Current behaviour is "merging" them
func FromArguments(args Arguments) (configuration Configuration, err error) {
	conf, err := createConfiguration(args.ConfigurationFile)
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

	configuration = Configuration{
		Action:            args.Action,
		BricksNames:       args.BricksNames,
		BricksSpecifiers:  extools.Deduplicate(append(conf.BricksSpecifiers, args.BricksSpecifiers...)),
		ConfigurationFile: args.ConfigurationFile,
		Format:            args.Format,
		Interactive:       !args.NonInteractive,
		Modules:           modules,
		Rooms:             rooms,
		OtherOptions:      extools.Deduplicate(append(conf.OtherOptions, args.OtherOptions...)),
	}

	return
}
