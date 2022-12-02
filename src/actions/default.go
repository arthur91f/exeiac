package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func DefaultBehaviour(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {

	statusCode = 3
	// a test just to use the interface arguments
	if infra != nil && args != nil {
		err = exargs.ErrBadArg{Reason: fmt.Sprintf(
			"Error: %s action not code yet", args.Action)}
	} else {
		err = exargs.ErrBadArg{Reason: "Error: infra and args are not setted"}
	}
	return
}

func Init(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {

	statusCode, err = DefaultBehaviour(infra, args, bricksToExecute)
	return
}

func ValidateCode(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {

	statusCode, err = DefaultBehaviour(infra, args, bricksToExecute)
	return
}
