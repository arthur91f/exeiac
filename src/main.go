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

var arguments exargs.Arguments

// The `init()` function runs before everything else.
// c.f. https://go.dev/doc/effective_go#init
func init() {
	flag.StringSliceVarP(&arguments.BricksSpecifiers, "bricks-specifiers", "s", []string{"selected"},
		fmt.Sprintf("A list of comma separated specifiers. Includes: %v", exargs.AvailableBricksSpecifiers))

	flag.StringVarP(&arguments.ConfigurationFile, "configuration-file", "c", "/etc/exeiac/conf.yml",
		"A path the a valid configuration file")

	flag.BoolVarP(&arguments.NonInteractive, "non-interactive", "I", false,
		"Allows for exeiac to run without user input")

	flag.StringVarP(&arguments.Format, "format", "f", "all",
		fmt.Sprintf("Define the format of the output. It matches the brick's specifiers values: %v", exargs.AvailableBricksSpecifiers))

	defaultModules := make(map[string]string)
	flag.StringToStringVarP(&arguments.Modules, "modules", "m", defaultModules,
		"A set of key/value pairs with the key being the name of the module, and the value its absolute path. This is useful if you want to setup a module on the command line. It uses the configuration file's modules declaration otherwise.")

	defaultRooms := make(map[string]string)
	flag.StringToStringVarP(&arguments.Rooms, "rooms", "r", defaultRooms,
		"A set of key/value pairs with the key being the name of the room, and the value its absolute path. This is useful if you want to setup a room on the command line. It uses the configuration file's rooms declaration otherwise.")

	flag.StringSliceVarP(&arguments.OtherOptions, "other-options", "o", []string{},
		`A list of options to pass straight to the module. Useful to execute a module with an action not handled by exeiac, or to provide extra-options to it. Flag with arguments need to be enclosed in double quotes (e.g. -o "--myflag myargument",-b)`)
}

func main() {
	var statusCode int

	flag.Parse()

	// The only remaining arguments are not flags. They match the action and the brickNames
	nonFlagArgs := flag.Args()
	arguments.Action, arguments.BricksNames = nonFlagArgs[0], nonFlagArgs[1:]

	configuration, err := exargs.FromArguments(arguments)
	if err != nil {
		fmt.Printf("%v\n> Error636a4c9e:main/main: unable to get arguments\n",
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

	// valid arguments (arg.brickNames are in infra.Bricks...)
	err = infra.ValidateConfiguration(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	// enrich bricks that we will execute
	infra.EnrichBricks()

	// get bricks selected
	var bricks exinfra.Bricks
	bricks, err = infra.GetBricksFromNames(arguments.BricksNames)
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
	// if arguments.action is in the list do that else use otherAction
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
