package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
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
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for plan action"}

		return exstatuscode.INIT_ERROR, err
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
			return exstatuscode.ENRICH_ERROR, err
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
		statusCode = exstatuscode.INIT_ERROR
		err = exinfra.ErrBadArg{Reason: fmt.Sprintf(
			"Error: format not valid for show action: %s", conf.Format)}

		return
	}

	return
}
