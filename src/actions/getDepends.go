package actions

import (
	"fmt"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
)

type depends struct {
	brick    *exinfra.Brick
	jsonPath string
	varName  string
}

func GetDepends(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(conf.BricksNames) != 1 {
		return exstatuscode.INIT_ERROR,
			exinfra.ErrBadArg{Reason: fmt.Sprintf("Error: you should specify a single brick for %s action", conf.Action)}
	}
	var bricks exinfra.Bricks
	var depends exinfra.Bricks

	bricks, err = infra.GetBricksFromNames(conf.BricksNames)
	if err != nil {
		statusCode = exstatuscode.RUN_ERROR

		return // TODO(arthur91f): check display
	}

	depends, err = infra.GetBricksThatCallthisOutput(bricks[0], conf.JsonPath)
	if err != nil {
		statusCode = exstatuscode.RUN_ERROR

		return // TODO(arthur91f): check display
	}

	sort.Sort(depends)

	switch conf.Format {
	case "name", "n":
		for _, b := range depends {
			fmt.Println(b.Name)
		}
	case "path", "p":
		for _, b := range depends {
			fmt.Println(b.Path)
		}
	case "all", "a": // brick:var_name <- jsonPath
		for _, b := range depends {
			inputs := b.GetInputsThatCallthisOutput(bricks[0], conf.JsonPath)
			for _, i := range inputs {
				fmt.Printf(
					"%s:%s = %s\n",
					b.Name, i.VarName, i.Dependency.From.JsonPath,
				)
			}
		}
	default:
		return exstatuscode.INIT_ERROR,
			exinfra.ErrBadArg{Reason: fmt.Sprintf(
				"Error: format \"%s\" not available for get-depends", conf.Format)}
	}

	return

}
