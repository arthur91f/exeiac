package main

import (
    "fmt"
    exeiacArgs "src/exeiac/arguments"
)

func main() {
    args, has_success := exeiacArgs.Get_arguments()
    if has_success {
        fmt.Println(args.Action)
        fmt.Println(args.Bricks_paths)
    } else {
        fmt.Println("Error: bad arguments")
    }
}

