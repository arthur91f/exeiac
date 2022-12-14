package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Show(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	// a test just to use the interface arguments
	if infra == nil && conf == nil {
		statusCode = 3
		err = exinfra.ErrBadArg{Reason: "Error: infra and args are not setted"}
		return
	}

	switch conf.Format {
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
		err = enrichDatas(bricksToExecute, infra)
		if err != nil {
			return 3, err
		}
		if len(bricksToExecute) == 1 {
			fmt.Println(string(bricksToExecute[0].Output))
		} else {
			for _, brick := range bricksToExecute {
				extools.DisplaySeparator(brick.Name)
				fmt.Println(string(brick.Output))
			}
		}
	default:
		statusCode = 3
		err = fmt.Errorf("Error: format not valid for show action: %s", conf.Format)
		return
	}
	return
}
