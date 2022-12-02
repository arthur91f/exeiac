package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

func DebugArgs(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute []string) (int, error) {
	fmt.Println(args)
	return 0, nil
}

func DebugInfra(infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute []string) (statusCode int, err error) {
	statusCode = 0
	fmt.Println(*infra)
	return
}
