package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
	"strings"
)

// Triggers a module execution for very single brick in `bricksToExecute`
// Ignores errors and calls the action in `args.Action` for every single brick,
// then prints out a summary of it all.
// Exit code matches 3 if an error occured, 0 otherwise.

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

	var linkedBricks exinfra.Bricks
	var depends exinfra.Bricks

	var bricks exinfra.Bricks
	bricks, err = infra.GetBricksFromNames(conf.BricksNames)
	if err != nil {
		statusCode = exstatuscode.RUN_ERROR

		return // TODO(arthur91f): check display
	}

	linkedBricks, err = infra.GetDirectNext(bricks[0])

	for _, b := range linkedBricks {
		for _, i := range b.Inputs {
			if i.Brick == bricks[0] {
				if areJsonPathsLinked(conf.JsonPath, i.JsonPath) {
					depends = append(depends, b)
					break
				}
			}
		}
	}

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
		return exstatuscode.INIT_ERROR,
			exinfra.ErrBadArg{Reason: fmt.Sprintf(
				"Error: format \"%s\" not available for get-depends", conf.Format)}
	default:
		return exstatuscode.INIT_ERROR,
			exinfra.ErrBadArg{Reason: fmt.Sprintf(
				"Error: format \"%s\" not available for get-depends", conf.Format)}
	}

	return

}

func areJsonPathsLinked(jsonpath1 string, jsonpath2 string) bool {

	// at this step assume that both are jsonpath

	jsonpathSanitized1 := jsonpath1
	if strings.HasSuffix(jsonpath1, ".*") {
		jsonpathSanitized1 = jsonpath1[0 : len(jsonpath1)-2]
	}

	jsonpathSanitized2 := jsonpath2
	if strings.HasSuffix(jsonpath2, ".*") {
		jsonpathSanitized2 = jsonpath2[0 : len(jsonpath2)-2]
	}

	jsonpathSlice1 := strings.Split(jsonpathSanitized1, ".")
	jsonpathSlice2 := strings.Split(jsonpathSanitized2, ".")

	var smallJsonpathSlice []string
	var longJsonpathSlice []string
	if len(jsonpathSlice1) > len(jsonpathSlice2) {
		longJsonpathSlice = jsonpathSlice1
		smallJsonpathSlice = jsonpathSlice2
	} else {
		longJsonpathSlice = jsonpathSlice2
		smallJsonpathSlice = jsonpathSlice1
	}

	for i, s := range smallJsonpathSlice {

		if s != longJsonpathSlice[i] {

			return false
		}

	}

	return true

}

// getOutput from/brick/name
// {
// 	database: {
// 		credentials: {
// 			username: toto
// 			password: tataTttt
// 		}
// 	}
// 	app: {
// 		...
// 	}
// }

// GetDepends from/brick/name:$.database.credentials.*

// i/brick1 creds from/brick/name:$.database.credentials
// i/brick2 creds from/brick/name:$.database.credentials.username
// i/brick2 creds from/brick/name:$.database.credentials.whatever
// i/brick3 mysql from/brick/name:$.database
// i/brick4 bricko from/brick/name:$.*

// $.*

// don't depends
// i/brick5 app from/brick/name:$.app
// i/brick6 we  from/brick/name:$.whatever
