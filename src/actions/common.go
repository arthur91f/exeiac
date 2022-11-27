package actions

import (
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

type Action struct {
	Name         string
	Behaviour    string
	UserCall     []string
	ModuleAction bool
}

var behaviourMap = map[string]func(*exinfra.Infra, *exargs.Arguments, []string) (int, error){
	"cd":            ChangeDirectory,
	"clean":         Clean,
	"help":          Help,
	"init":          Init,
	"lay":           Lay,
	"plan":          Plan,
	"remove":        Remove,
	"show":          Show,
	"validate_code": ValidateCode,
	"debug_args":    DebugArgs,
	"debug_infra":   DebugInfra,
	"default":       DefaultBehaviour,
}

func (a Action) Execute(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute []string) (statusCode int, err error) {

	statusCode, err = behaviourMap[a.Behaviour](infra, args, bricksToExecute)
	return
}

func CreateActionsMap(infra *exinfra.Infra, bricks []string) (
	m map[string]Action, err error) {

	m = map[string]Action{
		"cd": {
			Behaviour:    "cd",
			UserCall:     []string{"change_directory", "cd"},
			ModuleAction: false,
		},
		"clean": {
			Behaviour:    "clean",
			UserCall:     []string{"clean"},
			ModuleAction: true,
		},
		"help": {
			Behaviour:    "help",
			UserCall:     []string{"help"},
			ModuleAction: true,
		},
		"init": {
			Behaviour:    "init",
			UserCall:     []string{"init"},
			ModuleAction: true,
		},
		"lay": {
			Behaviour:    "lay",
			UserCall:     []string{"lay"},
			ModuleAction: true,
		},
		"plan": {
			Behaviour:    "plan",
			UserCall:     []string{"plan"},
			ModuleAction: true,
		},
		"remove": {
			Behaviour:    "remove",
			UserCall:     []string{"remove", "rm"},
			ModuleAction: true,
		},
		"show": {
			Behaviour:    "show",
			UserCall:     []string{"show"},
			ModuleAction: true,
		},
		"validate_code": {
			Behaviour:    "validate_code",
			UserCall:     []string{"validate_code"},
			ModuleAction: true,
		},
		"debug_args": {
			Behaviour:    "debug_args",
			UserCall:     []string{"debug_args"},
			ModuleAction: false,
		},
		"debug_infra": {
			Behaviour:    "debug_infra",
			UserCall:     []string{"debug_infra"},
			ModuleAction: false,
		},
	}
	for _, brick := range bricks {
		for _, action := range infra.Bricks[brick].Module.Actions {
			if _, exists := m[action]; !exists {
				m[action] = Action{
					Behaviour:    "default",
					UserCall:     []string{action},
					ModuleAction: true,
				}
			}
		}
	}
	return
}
