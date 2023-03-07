package main

import (
	"fmt"
	"os"
	exaction "src/exeiac/actions"
	exargs "src/exeiac/arguments"
	excompletion "src/exeiac/completion"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"

	flag "github.com/spf13/pflag"
)

// The `init()` function runs before everything else.
// c.f. https://go.dev/doc/effective_go#init
func init() {
	flag.Parse()
}

func main() {
	var statusCode int

	if exargs.Args.ShowUsage {
		flag.Usage()
		return
	}

	// The only remaining arguments are not flags. They match the action and the brickNames
	nonFlagArgs := flag.Args()

	// NOTE(half-shell): We need the configuration created to list all bricks so we bypass
	// the check made on the action and bricks if the "list bricks" flag is provided.
	if !exargs.Args.ListBricks {
		if len(nonFlagArgs) == 0 {
			fmt.Fprintln(os.Stderr, "argument missing: you need at least to specify one action")

			os.Exit(exstatuscode.INIT_ERROR)
		}

		exargs.Args.Action, exargs.Args.BricksNames = nonFlagArgs[0], nonFlagArgs[1:]
	}

	configuration, err := exargs.FromArguments(exargs.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(exstatuscode.INIT_ERROR)
	}

	if exargs.Args.ListBricks {
		excompletion.ListBricks(configuration)

		os.Exit(0)
	}

	// build infra representation
	infra, err := exinfra.CreateInfra(configuration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(exstatuscode.INIT_ERROR)
	}

	err = infra.ValidateConfiguration(&configuration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(exstatuscode.INIT_ERROR)
	}

	// enrich bricks that we will execute
	infra.EnrichBricks()

	// get bricks selected
	var bricks exinfra.Bricks
	bricks, err = infra.GetBricksFromNames(exargs.Args.BricksNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(exstatuscode.INIT_ERROR)
	}

	// get bricks specified by parameters
	var bricksToExecute exinfra.Bricks
	bricksToExecute, err = infra.GetCorrespondingBricks(bricks, configuration.BricksSpecifiers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(exstatuscode.INIT_ERROR)
	}
	// executeAction
	// if exargs.Args.action is in the list do that else use otherAction
	if behaviour, ok := exaction.BehaviourMap[configuration.Action]; ok {
		statusCode, err = behaviour(&infra, &configuration, bricksToExecute)
	} else {
		statusCode, err = exaction.BehaviourMap["default"](&infra, &configuration, bricksToExecute)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	exaction.Test129_4(&infra)

	os.Exit(statusCode)
}
