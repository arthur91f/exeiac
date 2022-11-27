package executionFlow

import (
	"fmt"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

type executionStep struct {
	Action     string
	Brick      *exinfra.Brick
	State      string
	StatusCode int
}

type ExecutionPlan []executionStep

func (e executionStep) String() string {
	var s string
	a := e.Action
	b := e.Brick.Name
	if e.State != "" {
		if e.StatusCode > 0 {
			s = fmt.Sprintf("%s(%d)", e.State, e.StatusCode)
		} else { // success or abnormal failure without return code
			s = e.State
		}
	} else {
		s = "todo"
	}
	return fmt.Sprintf("%s:%s: %s", a, b, s)
}

func (e ExecutionPlan) String() string {
	var str string
	for _, step := range e {
		str = fmt.Sprintf("%s%s\n", str, step.String())
	}
	return str
}

func (e ExecutionPlan) SkipFromIndex(index int) string {
	for i := index + 1; index < len(e); i++ {
		e[i].State = "skipped"
	}
	return e.String()
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
