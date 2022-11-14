package actions

import (
	"fmt"
	exinfra "src/exeiac/infra"
)

func ShowInfra(i exinfra.Infra) {
	fmt.Println("Infra:")

	fmt.Println("  Modules:")
	for _, module := range i.Modules {
		fmt.Println("  - name: " + module.Name)
		fmt.Println("    path: " + module.Path)
	}

	fmt.Println("  Rooms:")
	for _, room := range i.Rooms {
		fmt.Println("  - name: " + room.Name)
		fmt.Println("    path: " + room.Path)
	}

	fmt.Println("  Bricks:")
	for _, brick := range i.Bricks {
		fmt.Println("  - name: " + brick.Name)
		fmt.Println("    path: " + brick.Path)
		fmt.Printf("    isElementary: %t\n", brick.IsElementary)
		fmt.Println("    confFile: " + brick.ConfigurationFilePath)
	}
}
