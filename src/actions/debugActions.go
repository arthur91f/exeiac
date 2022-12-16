package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func DebugArgs(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	fmt.Println(conf)

	return
}

func DebugInfra(
	infra *exinfra.Infra, conf *exargs.Configuration, bricksToExecute exinfra.Bricks) (
	statusCode int, err error) {
	statusCode = 0
	fmt.Println(infra)
	fmt.Printf("bricksToExecute:")
	if len(bricksToExecute) == 0 {
		fmt.Println(" []")
	} else {
		for _, b := range bricksToExecute {
			fmt.Printf("\n  - %d:%s", b.Index, b.Name)
		}
		fmt.Println("")
	}

	return
}
