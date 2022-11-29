package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	exaction "src/exeiac/actions"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

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

	// enrich bricks that we will execute
	enrichBricks(&infra)

	// TODO(arthur91f): func getBricksToExecute(args.BricksNames args.Specifier)
	var bricksToExecute []string
	bricksToExecute = args.BricksNames

	// executeAction
	// if args.action is in the list do that else use otherAction
	if behaviour, ok := exaction.BehaviourMap[args.Action]; ok {
		statusCode, err = behaviour(&infra, &args, bricksToExecute)
	} else {
		statusCode, err = exaction.BehaviourMap["default"](&infra, &args, bricksToExecute)
	}
	if err != nil {
		log.Printf("%v\n", err)
	}
	os.Exit(statusCode)
}

var availableBricksSpecifiers = []string{
	"linked_previous", "all_previous", "lp", "ap",
	"direct_previous", "dp",
	"selected", "s",
	"direct_next", "dn",
	"linked_next", "all_next", "ln", "an"}

func validArgBricksAreInInfra(infra *exinfra.Infra, args *exargs.Arguments) error {
	// valid that args.BricksNames items are valid
	for _, arg := range args.BricksNames {
		if _, ok := infra.Bricks[arg]; !ok {
			return exargs.ErrBadArg{Reason: "Brick doesn't exist:", Value: arg}
		}
	}

	// valid BricksSpecifiers
	for _, specifier := range args.BricksSpecifiers {
		if !extools.ContainsString(availableBricksSpecifiers, specifier) {
			return exargs.ErrBadArg{Reason: "Brick's specifier doesn't exist:",
				Value: specifier}
		}
	}
	return nil
}

func enrichBricks(infra *exinfra.Infra) {
	// TODO(arthur91f): remove log.Fatal
	for _, b := range infra.Bricks {
		if b.IsElementary {
			conf, err := exinfra.BrickConfYaml{}.New(b.ConfigurationFilePath)
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}

			err = b.Enrich(conf, infra)
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}
			err = b.Module.LoadAvailableActions()
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}
		}
	}
}

func getBricksToExecute(
	infra *exinfra.Infra, args *exargs.Arguments) (
	bricksToExecute []string, err error) {

	var bricks exinfra.Bricks
	for _, brickName := range args.BricksNames {
		bricks = append(bricks, infra.Bricks[brickName])
	}

	var bricksSpecified exinfra.Bricks
	var bricksToAdd exinfra.Bricks
	for _, specifier := range args.BricksSpecifiers {
		bricksToAdd = exinfra.Bricks{}
		switch specifier {
		case "linked_previous", "all_previous", "lp", "ap":
			for _, brick := range bricks {
				bs, err := infra.GetLinkedPrevious(brick.Name) // recursive GetDirectPrevious
				if err != nil {
					// Handle error somehow
					log.Fatal(err)
				}
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "direct_previous", "dp":
			for _, brick := range bricks {
				bs, err := infra.GetDirectPrevious(brick.Name)
				if err != nil {
					// Handle error somehow
					log.Fatal(err)
				}
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "selected", "s":
			for _, brick := range bricks {
				bs, err := infra.GetSubBricks(brick.Name) // done
				if err != nil {
					// Handle error somehow
					log.Fatal(err)
				}
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "direct_next", "dn":
			for _, brick := range bricks {
				bs, err := infra.GetDirectNext(brick.Name)
				if err != nil {
					// Handle error somehow
					log.Fatal(err)
				}
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "linked_next", "all_next", "ln", "an":
			for _, brick := range bricks {
				var bs exinfra.Bricks
				err := infra.GetLinkedNext(&bs, brick.Name)
				if err != nil {
					// Handle error somehow
					log.Fatal(err)
				}
				bricksToAdd = append(bricksToAdd, bs...)
			}
		default:
			err = exargs.ErrBadArg{Reason: "Brick's specifier doesn't exist:",
				Value: specifier}
			return
		}

		bricksSpecified = append(bricksSpecified, bricksToAdd...)
	}

	sort.Sort(bricks)
	bricks.RemoveDuplicates()

	return
}
