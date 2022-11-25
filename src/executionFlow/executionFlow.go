package executionFlow

import (
	"fmt"
	"os/exec"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

type executionStep struct {
	Action    string
	Brick     *exinfra.Brick
	State     string
	ExitError *exec.ExitError
}

type ExecutionPlan []executionStep

func (e ExecutionPlan) PrintPlan() {
	fmt.Println("Here the execution plan")
	for _, step := range e {
		fmt.Printf("- action: %s\n", step.Action)
		fmt.Printf("  brick: %s\n", step.Brick.Name)
	}
}

func CreateExecutionPlan(i *exinfra.Infra, action string, brickNames []string) (ExecutionPlan, error) {
	var inputBricks exinfra.Bricks
	var bricks exinfra.Bricks
	var executionPlan ExecutionPlan

	for _, n := range brickNames {
		if b, ok := i.Bricks[n]; ok {
			inputBricks = append(inputBricks, b)
		} else {
			return executionPlan, exargs.ErrBadArg{
				Reason: "brick not recognized: ", Value: n}
		}
	}

	// TODO: modify bricks with specifiers

	// We break super-bricks down to elementary bricks
	for _, b := range inputBricks {
		if i.Bricks[b.Name].IsElementary {
			bricks = append(bricks, b)
		} else {
			subBricks := i.GetSubBricks(b)
			bricks = append(bricks, subBricks...)
		}
	}

	sort.Sort(bricks)

	// create step
	for _, b := range bricks {
		executionPlan = append(executionPlan, executionStep{
			Action: action,
			Brick:  b,
		})
	}

	return executionPlan, nil
}
