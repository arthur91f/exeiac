package actions

import (
	"fmt"
	"sort"
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

func Test129_2_1(infra *exinfra.Infra) {
	var bricksToDisplay exinfra.Bricks
	for n, b := range infra.Bricks {
		bricksToDisplay = append(bricksToDisplay, b)
		if n != b.Name {
			fmt.Printf("ERR129_2.1:1: %s != %s\n", n, b.Name)
		}
	}
	sort.Sort(bricksToDisplay)
	for i, b := range bricksToDisplay {
		fmt.Printf("%d:%d:%t:%s\n", i, b.Index, b.IsElementary, b.Name)
	}
}
