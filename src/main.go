package main

import (
    "fmt"
    exeiacArgs "src/exeiac/arguments"
    "os"
)

func action_display_args(args exeiacArgs.Arguments) {
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

    fmt.Printf("  interactive: %t", args.Interactive)

    if len(args.ModulesList) == 0 {
        fmt.Println("  modules_list: []")
    } else {
        fmt.Println("  modules_list:")
        for _, module := range args.ModulesList {
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

    if len(args.RoomsList) == 0 {
        fmt.Println("  rooms_list: []")
    } else {
        fmt.Println("  rooms_list:")
        for _, room := range args.RoomsList {
            fmt.Println("  - name:" + room.Name)
            fmt.Println("    path:" + room.Path)
        }
    }
}

func main() {
    args, err := exeiacArgs.GetArguments()
    if err != nil {
        fmt.Printf("Error:exeiac: %w\n", err)
        os.Exit(1)
    }
    action_display_args(args)
}

