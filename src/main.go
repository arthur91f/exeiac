package main

import (
	"fmt"
	"log"
	"os"
	exactions "src/exeiac/actions"
	exargs "src/exeiac/arguments"
	exexec "src/exeiac/executionFlow"
	exinfra "src/exeiac/infra"
)

func main() {
	// get arguments
	args, err := exargs.GetArguments()
	if err != nil {
		fmt.Printf("%v\n> Error636a4c9e:main/main: unable to get arguments\n",
			err)
		os.Exit(1)
	}
	exactions.ShowArgs(args)

	// build infra representation
	infra, err := exinfra.CreateInfra(args.Rooms, args.Modules)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}
	fmt.Println(infra)

	// build executionPlan
	// TODO: Replace the last arguments to contain a list of brick names
	// executionPlan, err := exexec.ExecutionPlan{}.New(infra, args.Action, args.brickNames)
	executionPlan, err := exexec.CreateExecutionPlan(&infra, args.Action, []string{args.Brick})
	if err != nil {
		fmt.Printf("%v\n> Error6373c57e:main/main: "+
			"unable to get the executionPlan\n", err)
		os.Exit(1)
	}

	for _, step := range executionPlan {
		conf, err := exinfra.BrickConfYaml{}.New(step.Brick.ConfigurationFilePath)
		if err != nil {
			log.Fatalf("Unable to load brick's configuration file: %s", err)
		}

		err = step.Brick.Enrich(conf, &infra)
		if err != nil {
			log.Fatalf("Unable to enrich brick: %s", err)
		}

		err = step.Brick.Module.LoadAvailableActions()
		if err != nil {
			log.Fatalf("Unable to load action for module %s: %s", step.Brick, err)
		}
	}

	executionPlan.PrintPlan()
}
