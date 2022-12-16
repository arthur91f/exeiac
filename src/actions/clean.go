package actions

import (
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func Clean(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	statusCode = 3
	// a test just to use the interface arguments
	if infra != nil && conf != nil {
		err = exinfra.ErrBadArg{Reason: "Error: clean action not code yet"}
	} else {
		err = exinfra.ErrBadArg{Reason: "Error: infra and args are not setted"}
	}
	return
}
