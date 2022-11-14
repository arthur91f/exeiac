package main

import (
	"fmt"
	"os"
	exactions "src/exeiac/actions"
	exargs "src/exeiac/arguments"
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
	infra, err := exinfra.Infra{}.New(args.RoomsList, args.ModulesList)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}

	infra.Display()
}
