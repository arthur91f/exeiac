package actions

import (
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func Show(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute []string) (statusCode int, err error) {

	statusCode = 3
	// a test just to use the interface arguments
	if infra != nil && args != nil {
		err = exargs.ErrBadArg{Reason: "Error: show action not code yet"}
	} else {
		err = exargs.ErrBadArg{Reason: "Error: infra and args are not setted"}
	}
	return
}
