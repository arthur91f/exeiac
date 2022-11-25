package actions

import (
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

var actionsMap = map[string]func(*exinfra.Infra, *exargs.Arguments) (int, error){
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
	// the personnal actions not implemented
}

func validArgs(infra *exinfra.Infra, args *exargs.Arguments) error {
	// valid that args.BricksNames items are valid
	for _, arg := range args.BricksNames {
		if _, ok := infra.Bricks[arg]; !ok {
			return exargs.ErrBadArg{Reason: "Brick doesn't exist:", Value: arg}
		}
	}
	return nil
}

func Execute(infra *exinfra.Infra, args *exargs.Arguments) (statusCode int, err error) {
	err = validArgs(infra, args)
	if err != nil {
		statusCode = 3
		return
	}
	statusCode, err = actionsMap[args.Action](infra, args)
	return
}
