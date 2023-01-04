package arguments

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

var Args Arguments

func init() {
	flag.StringSliceVarP(&Args.BricksSpecifiers, "bricks-specifiers", "s", []string{"selected"},
		fmt.Sprintf(`A list of comma separated specifiers.
Includes: %v`, AvailableBricksSpecifiers))

	flag.StringVarP(&Args.ConfigurationFile, "configuration-file", "c", "",
		"A path the a valid configuration file")

	flag.BoolVarP(&Args.NonInteractive, "non-interactive", "I", false,
		"Allows for exeiac to run without user input")

	flag.BoolVarP(&Args.Interactive, "interactive", "i", false,
		"Forces exeiac to wait for user input when needed. Overwrites the --non-interactive (-I) flag.")

	flag.StringVarP(&Args.Format, "format", "f", "all",
		fmt.Sprintf(`Define the format of the output. It matches the brick's specifiers values
Includes: %v`, AvailableBricksFormat))

	defaultModules := make(map[string]string)
	flag.StringToStringVarP(&Args.Modules, "modules", "m", defaultModules,
		`A set of key/value pairs with the key being the name of the module,
and the value its absolute path. This is useful if you want to setup
a module on the command line. It uses the configuration file's modules
declaration otherwise.`)

	defaultRooms := make(map[string]string)
	flag.StringToStringVarP(&Args.Rooms, "rooms", "r", defaultRooms,
		`A set of key/value pairs with the key being the name of the room,
and the value its absolute path. This is useful if you want to setup
a room on the command line. It uses the configuration file's rooms
declaration otherwise.`)

	flag.StringSliceVarP(&Args.OtherOptions, "other-options", "o", []string{},
		`A list of options to pass straight to the module. Useful to execute a
module with an action not handled by exeiac, or to provide extra-options
to it. Flag with arguments need to be enclosed in double quotes
(e.g. -o "--myflag myargument",-b)`)

	flag.BoolVarP(&Args.ShowUsage, "help", "h", false, "Show exeiac's help")

	flag.BoolVarP(&Args.ListBricks, "list-bricks", "l", false, "List all the bricks from all rooms")

	flag.Usage = func() {
		fmt.Println("Usage: exeiac ACTION (BRICK_PATH|BRICK_NAME) [OPTION...]")
		fmt.Println()
		fmt.Println("ACTION:")
		fmt.Println("  init: get some dependencies, typically download terraform modules or ansible deps")
		fmt.Println("  plan: a dry run to check what we want to lay")
		fmt.Println("  lay: lay the brick on the wall. Run the IaC with the right tools")
		fmt.Println("  remove: remove a brick from your wall to destroy it properly.")
		fmt.Println("  validate_code: validate if the syntaxe is ok")
		fmt.Println("  help: display this help or the specified help for the brick")
		fmt.Println("  show: display brick attributes (depends of the format option choosen)")
		fmt.Println("  clean: remove all files created by exeiac")
		fmt.Println()
		fmt.Println("OPTION:")
		flag.PrintDefaults()
	}
}
