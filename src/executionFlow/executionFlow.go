package executionFlow

import (
	"fmt"
	"os/exec"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
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

func (e ExecutionPlan) New(i *exinfra.Infra, args *exargs.Arguments) (ExecutionPlan, error) {

	action := args.Action
	var bricksIndexes []int
	var executionPlan ExecutionPlan

	// arg.Brick
	if args.Brick != "" {
		if ibrick, err := i.GetBrickIndexWithName(args.Brick); err == nil {
			bricksIndexes = append(bricksIndexes, ibrick)
		} else if ibrick, err := i.GetBrickIndexWithPath(args.Brick); err == nil {
			bricksIndexes = append(bricksIndexes, ibrick)
		} else {
			return executionPlan, exargs.ErrBadArg{
				Reason: "brick not recognized: ", Value: args.Brick}
		}
	}

	// arg.BricksPaths
	if len(args.BricksPaths) != 0 {
		for _, brickPath := range args.BricksPaths {
			ibrick, err := i.GetBrickIndexWithPath(brickPath)
			if err != nil {
				return executionPlan, exargs.ErrBadArg{
					Reason: "brick not recognized: ", Value: brickPath}
			}
			bricksIndexes = append(bricksIndexes, ibrick)
		}
	}

	// arg.BricksNames
	if len(args.BricksNames) != 0 {
		for _, brickName := range args.BricksNames {
			ibrick, err := i.GetBrickIndexWithName(brickName)
			if err != nil {
				return executionPlan, exargs.ErrBadArg{
					Reason: "brick not recognized: ", Value: brickName}
			}
			bricksIndexes = append(bricksIndexes, ibrick)
		}
	}

	// TODO: modify bricks with specifiers

	// GetElementaryBricks
	var tempIndexes []int
	for _, index := range bricksIndexes {
		if i.Bricks[index].IsElementary {
			tempIndexes = append(tempIndexes, index)
		} else {
			tempIndexes = append(tempIndexes, i.GetSubBricksIndexes(index)...)
		}
	}
	bricksIndexes = tempIndexes

	// Remove duplicates and sort the list
	bricksIndexes = extools.GetIntSliceWithoutDuplicates(bricksIndexes)
	sort.Ints(bricksIndexes)

	// create step
	for _, brickIndex := range bricksIndexes {
		executionPlan = append(executionPlan, executionStep{
			Action: action,
			Brick:  &(i.Bricks[brickIndex]),
		})
	}
	return executionPlan, nil
}
