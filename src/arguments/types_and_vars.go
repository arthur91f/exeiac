package arguments

import extools "src/exeiac/tools"

type Arguments struct {
	Action           string
	BricksPaths      []string
	BricksSpecifiers []string
	Interactive      bool
	Modules          []extools.NamePathBinding
	OutputSpecifier  string
	OtherOptions     []string
	Rooms        []extools.NamePathBinding
}

func getDefaultArguments() Arguments {
	return Arguments{
		Action:           "",
		BricksPaths:      []string{},
		BricksSpecifiers: []string{"selected"},
		Interactive:      true,
		Modules:          []extools.NamePathBinding{},
		OutputSpecifier:  ".",
		OtherOptions:     []string{},
		Rooms:        []extools.NamePathBinding{},
	}
}

type exeiacConf struct {
	Rooms   []extools.NamePathBinding `yaml:"rooms"`
	Modules []extools.NamePathBinding `yaml:"modules"`
	DefaultArgs struct {
		NonInteractive   bool   `yaml:"non_interactive"`
		BricksSpecifiers string `yaml:"bricks_specifiers"`
		OtherOptions     string `yaml:"other_options"`
	} `yaml:"default_arguments"`
}

var actions_list = []string{
	"plan", "lay", "remove", "output", "init", "validate_code", "help",
	"show_input", "list_elementary_bricks", "cd",
	"get_brick_path", "get_brick_name"}
