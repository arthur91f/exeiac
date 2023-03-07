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

func Test129_4(infra *exinfra.Infra) {
	var bricksToDisplay exinfra.Bricks
	for n, b := range infra.Bricks {
		bricksToDisplay = append(bricksToDisplay, b)
		if b.EnrichError != nil {
			fmt.Printf("ERR129_4.3:1: %s:%v\n", n, b.EnrichError)
		}
	}
	sort.Sort(bricksToDisplay)
	for _, b := range bricksToDisplay {
		fmt.Print(s(b))
	}
}

func sb(b bool) string {
	if b {
		return "T"
	}
	return "F"
}

func sd(i int) string {
	return fmt.Sprintf("%d", i)
}

func s(b *exinfra.Brick) (s string) {
	// write sorted directPrevious
	dps := b.DirectPrevious
	sort.Sort(dps)
	dpsString := ""
	for _, dp := range dps {
		dpsString = fmt.Sprintf("%s%d-", dpsString, dp.Index)
	}

	// write sorted Inputs
	is := exinfra.CreateInputs(b.Inputs)
	is.SortInputByVarname()
	isString := ""
	for _, i := range is {
		isString = fmt.Sprintf("%s  %s\n", isString, si(i, b))
	}

	s = fmt.Sprintf("%d-%s:%s:%s\n%s",
		b.Index, b.Name,
		sb(b.IsElementary),
		dpsString, isString,
	)
	return
}

func si(i exinfra.Input, b *exinfra.Brick) (s string) {
	return fmt.Sprintf("%d:%s<%d:%s", b.Index, i.VarName, i.Brick.Index, i.JsonPath)
}
