package main

import (
	"fmt"
	"log"
	"os"
	//exactions "src/exeiac/actions"
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
	//exactions.ShowArgs(args)

	// build infra representation
	infra, err := exinfra.Infra{}.New(args.Rooms, args.Modules)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}
	fmt.Println(infra.String())

	// build executionPlan
	executionPlan, err := exexec.ExecutionPlan{}.New(infra, &args)
	if err != nil {
		fmt.Printf("%v\n> Error6373c57e:main/main: "+
			"unable to get the executionPlan\n", err)
		os.Exit(1)
	}

	for _, ep := range executionPlan {
		if ep.Brick.IsElementary {
			conf, err := exinfra.BrickConfYaml{}.New(ep.Brick.ConfigurationFilePath)
			if err != nil {
				log.Fatal("An issue happened while reading bricks' configurations\n", err)
			}

			ep.Brick.Enrich(conf, infra)
			ep.Brick.Module.LoadAvailableActions()
		}
	}

	executionPlan.PrintPlan()
}
