package actions

import (
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

var BehaviourMap = map[string]func(*exinfra.Infra, *exargs.Arguments, exinfra.Bricks) (int, error){
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
}
