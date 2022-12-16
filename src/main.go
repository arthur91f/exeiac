package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	exaction "src/exeiac/actions"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
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
	exargs.Args.Action, exargs.Args.BricksNames = nonFlagArgs[0], nonFlagArgs[1:]

	configuration, err := exargs.FromArguments(exargs.Args)
	if err != nil {
		fmt.Printf("%v\n> Error636a4c9e:main/main: unable to get exargs.Args\n",
			err)
		os.Exit(2)
	}

	// build infra representation
	infra, err := exinfra.CreateInfra(configuration)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}

	err = infra.ValidateConfiguration(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	// enrich bricks that we will execute
	infra.EnrichBricks()

	// get bricks selected
	var bricks exinfra.Bricks
	bricks, err = infra.GetBricksFromNames(exargs.Args.BricksNames)
	if err != nil {
		log.Fatal(err)
	}

	// get bricks specified by parameters
	var bricksToExecute exinfra.Bricks
	bricksToExecute, err = infra.GetCorrespondingBricks(bricks, configuration.BricksSpecifiers)
	if err != nil {
		log.Fatal(err)
	}

	// executeAction
	// if exargs.Args.action is in the list do that else use otherAction
	if behaviour, ok := exaction.BehaviourMap[configuration.Action]; ok {
		statusCode, err = behaviour(&infra, &configuration, bricksToExecute)
	} else {
		statusCode, err = exaction.BehaviourMap["default"](&infra, &configuration, bricksToExecute)
	}

	if err != nil {
		log.Printf("%v\n", err)
	}

	os.Exit(statusCode)
}
