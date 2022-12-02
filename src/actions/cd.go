package actions

import (
	"os"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func ChangeDirectory(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute []string) (statusCode int, err error) {

	if len(bricksToExecute) == 1 {
		err = os.Chdir(infra.Bricks[bricksToExecute[0]].Path)
	} else if len(args.BricksNames) == 0 {
		err = exargs.ErrBadArg{
			Reason: "You haven't specify any target brick to change directory"}
	} else {
		err = exargs.ErrBadArg{Reason: "To many bricks for cd command"}
	}
	if err != nil {
		statusCode = 3
	}
	return
}
