package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
)

func ShowArgs(args exargs.Arguments) {
	fmt.Println("Arguments:")
	fmt.Println("  action: " + args.Action)

	if len(args.BricksPaths) == 0 {
		fmt.Println("  bricks_paths: []")
	} else {
		fmt.Println("  bricks_paths:")
		for _, brick_path := range args.BricksPaths {
			fmt.Println("  - " + brick_path)
		}
	}

	if len(args.BricksSpecifiers) == 0 {
		fmt.Println("  bricks_specifiers: []")
	} else {
		fmt.Println("  bricks_specifiers:")
		for _, specifier := range args.BricksSpecifiers {
			fmt.Println("  - " + specifier)
		}
	}

	fmt.Printf("  interactive: %t\n", args.Interactive)

	if len(args.Modules) == 0 {
		fmt.Println("  modules_list: []")
	} else {
		fmt.Println("  modules_list:")
		for _, module := range args.Modules {
			fmt.Println("  - name:" + module.Name)
			fmt.Println("    path:" + module.Path)
		}
	}

	fmt.Println("  output_specifier: " + args.OutputSpecifier)

	if len(args.OtherOptions) == 0 {
		fmt.Println("  other_options: []")
	} else {
		fmt.Println("  other_options:")
		for _, opt := range args.OtherOptions {
			fmt.Println("  - " + opt)
		}
	}

	if len(args.Rooms) == 0 {
		fmt.Println("  rooms_list: []")
	} else {
		fmt.Println("  rooms_list:")
		for _, room := range args.Rooms {
			fmt.Println("  - name:" + room.Name)
			fmt.Println("    path:" + room.Path)
		}
	}
}
