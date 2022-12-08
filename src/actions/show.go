package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Show(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {

	// a test just to use the interface arguments
	if infra == nil && args == nil {
		statusCode = 3
		err = exargs.ErrBadArg{Reason: "Error: infra and args are not setted"}
		return
	}

	switch args.Format {
	case "path", "p":
		for _, brick := range bricksToExecute {
			fmt.Println(brick.Path)
		}
	case "name", "n":
		for _, brick := range bricksToExecute {
			fmt.Println(brick.Name)
		}
	case "all", "a":
		for _, brick := range bricksToExecute {
			fmt.Println(brick)
		}
	case "output", "outputs", "o":
		enrichDatas(bricksToExecute, infra)
		if len(bricksToExecute) == 1 {
			fmt.Println(bricksToExecute[0].Output)
		} else {
			for _, brick := range bricksToExecute {
				extools.DisplaySeparator(brick.Name)
				fmt.Println(brick.Output)
			}
		}
	default:
		statusCode = 3
		err = fmt.Errorf("Error: format not valid for show action: %s", args.Format)
		return
	}
	return
}
