package actions

import (
	"fmt"
	"sort"
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
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for show action"}

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
		var bricksToOutput exinfra.Bricks

		bricksToOutput, err = getBricksToOutput(bricksToExecute, infra, conf.Action)
		if err != nil {
			return exstatuscode.ENRICH_ERROR, err
		}

		bricksToOutput = append(bricksToOutput, bricksToExecute...)
		bricksToOutput = exinfra.RemoveDuplicates(bricksToOutput)
		sort.Sort(bricksToOutput)

		err = enrichOutputs(bricksToOutput)
		if err != nil {
			return exstatuscode.ENRICH_ERROR, err
		}

		for _, brick := range bricksToExecute {
			fmt.Println(brick)
		}
	case "output", "outputs", "o":
		var bricksToOutput exinfra.Bricks

		bricksToOutput, err = getBricksToOutput(bricksToExecute, infra, conf.Action)
		if err != nil {
			return exstatuscode.ENRICH_ERROR, err
		}

		bricksToOutput = append(bricksToOutput, bricksToExecute...)
		bricksToOutput = exinfra.RemoveDuplicates(bricksToOutput)
		sort.Sort(bricksToOutput)

		err = enrichOutputs(bricksToOutput)
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
