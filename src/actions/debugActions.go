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
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	fmt.Printf("Infra:\n%v\n", infra)
	fmt.Printf("bricksToExecute: [\n%v\n]", bricksToExecute)

	return
}
