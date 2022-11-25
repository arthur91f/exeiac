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

var actionsMap = map[string]func(*exinfra.Infra, *exargs.Arguments) (int, error){
	"cd":            ChangeDirectory,
	"clean":         Clean,
	"help":          Help,
	"init":          Init,
	"lay":           Lay,
	"plan":          Plan,
	"remove":        Remove,
	"show":          Show,
	"validate_code": ValidateCode,
	"debug_args":    DebugArgs,
	"debug_infra":   DebugInfra,
	// the personnal actions not implemented
}

func main() {
	var statusCode int
	// get arguments
	args, err := exargs.GetArguments()
	if err != nil {
		fmt.Printf("%v\n> Error636a4c9e:main/main: unable to get arguments\n",
			err)
		os.Exit(2)
	}

	// build infra representation
	infra, err := exinfra.CreateInfra(args.Rooms, args.Modules)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}

	// valid arguments (arg.brickNames are in infra.Bricks...)
	err = validArgBricksAreInInfra(&infra, &args)
	if err != nil {
		log.Fatal(err)
	}

	// TODO(arthur91f): func getBricksToExecute(args.BricksNames args.Specifier)
	var bricksToExecute []string
	bricksToExecute = args.BricksNames

	// enrich bricks that we will execute
	err = enrichBricks(&infra, bricksToExecute)
	if err != nil {
		log.Fatal(err)
	}

	// executeAction
	// if args.action is in the list do that else use otherAction
	statusCode, err = actionsMap[args.Action](&infra, &args, bricksToExecute)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	os.Exit(statusCode)

}

func validArgBricksAreInInfra(infra *exinfra.Infra, args *exargs.Arguments) error {
	// valid that args.BricksNames items are valid
	for _, arg := range args.BricksNames {
		if _, ok := infra.Bricks[arg]; !ok {
			return exargs.ErrBadArg{Reason: "Brick doesn't exist:", Value: arg}
		}
	}

	// TODO(arthur91f): valid the action
	//   if it's a known action -> ok
	//   if it's a not known action
	return nil
}

func enrichBricks(infra *exinfra.Infra, bricks []string) error {
	for _, b := range bricks {
		conf, err := exinfra.BrickConfYaml{}.New(infra.Bricks[b].ConfigurationFilePath)
		if err != nil {
			log.Fatalf("Unable to load brick's configuration file: %s", err)
		}

		err = infra.Bricks[b].Enrich(conf, infra)
		if err != nil {
			log.Fatalf("Unable to enrich brick: %s", err)
		}

		err = infra.Bricks[b].Module.LoadAvailableActions()
		if err != nil {
			log.Fatalf("Unable to load action for module %s: %s", b, err)
		}
	}
	return nil
}
