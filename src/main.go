package main

import (
	"fmt"
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
	infra, err := exinfra.Infra{}.New(args.Rooms, args.Modules)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}
	infra.Display()

	// build executionPlan
	executionPlan, err := exexec.ExecutionPlan{}.New(infra, &args)
	if err != nil {
		fmt.Printf("%v\n> Error6373c57e:main/main: "+
			"unable to get the executionPlan\n", err)
		os.Exit(1)
	}
	executionPlan.PrintPlan()
}
