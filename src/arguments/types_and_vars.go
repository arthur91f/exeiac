package arguments

type NamePathMapping struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func are_NamePathMappings_equal(m1 NamePathMapping, m2 NamePathMapping) (bool, bool) {
	are_names_equals := false
	are_path_equals := false

	if m1.Name == m2.Name {
		are_names_equals = true
	}
	if m1.Path == m2.Path {
		are_path_equals = true
	}
	return are_names_equals, are_path_equals
}

type Arguments struct {
	Action           string
	BricksPaths      []string
	BricksSpecifiers []string
	Interactive      bool
	ModulesList      []NamePathMapping
	OutputSpecifier  string
	OtherOptions     []string
	RoomsList        []NamePathMapping
}

func getDefaultArguments() Arguments {
	return Arguments{
		Action:           "",
		BricksPaths:      []string{},
		BricksSpecifiers: []string{"selected"},
		Interactive:      true,
		ModulesList:      []NamePathMapping{},
		OutputSpecifier:  ".",
		OtherOptions:     []string{},
		RoomsList:        []NamePathMapping{},
	}
}

type exeiacConf struct {
	RoomsList   []NamePathMapping `yaml:"rooms_list"`
	ModulesList []NamePathMapping `yaml:"modules_list"`
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
