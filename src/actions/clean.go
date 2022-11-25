package actions

import (
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func Clean(infra *exinfra.Infra, args *exargs.Arguments) (statusCode int, err error) {
	statusCode = 3
	// a test just to use the interface arguments
	if infra != nil && args != nil {
		err = exargs.ErrBadArg{Reason: "Error: clean action not code yet"}
	} else {
		err = exargs.ErrBadArg{Reason: "Error: infra and args are not setted"}
	}
	return
}
