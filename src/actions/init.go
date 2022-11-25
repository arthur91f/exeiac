package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exexec "src/exeiac/executionFlow"
	exinfra "src/exeiac/infra"
)

func Init(infra *exinfra.Infra, args *exargs.Arguments) (statusCode int, err error) {
	statusCode = 0
	var executionPlan exexec.ExecutionPlan
	var exitCode int

	executionPlan, err = exexec.CreateExecutionPlan(
		infra, args.Action, args.BricksNames)
	if err != nil {
		statusCode = 3
		return
	}

	for i, step := range executionPlan {
		fmt.Printf("-- %s --", step.Brick.Name)
		exitCode, err = step.Brick.ExecuteInteractive(step.Action)
		if err != nil {
			statusCode = 103
			executionPlan[i].State = fmt.Sprintf("Error")
			executionPlan[i].StatusCode = -1
			str := executionPlan.SkipFromIndex(i)
			fmt.Printf("%s", str)
			return
		} else if exitCode == 0 {
			executionPlan[i].State = "success"
			executionPlan[i].StatusCode = exitCode
		} else {
			statusCode = 3
			executionPlan[i].State = "failed"
		}
		fmt.Println("")
	}

	fmt.Printf("%s", executionPlan.String())
	return
}
