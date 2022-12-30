package actions

import (
	"fmt"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Remove(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for \"remove\" action"}

		return 3, err
	}

	if conf.Interactive {
		fmt.Println("Here, the bricks list to remove :")
		fmt.Print(bricksToExecute)

		// NOTE(half-shell): We might change this behavior to only ask for a "\n" input
		// instead of a Y/N choice.
		confirm, err := extools.AskConfirmation("\nDo you want to continue ?")

		if err != nil {
			return 3, err
		} else if !confirm {
			return 0, nil
		}
	}

	err = enrichDatas(bricksToExecute, infra)
	if err != nil {
		return 3, err
	}

	skipFollowing := false
	execSummary := make(ExecSummary, len(bricksToExecute))

	sort.Sort(sort.Reverse(bricksToExecute))

	for i, b := range bricksToExecute {

		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		// skip if an error was encounter before
		if skipFollowing {
			report.Status = "SKIP"
			execSummary[i] = report
			continue
		}

		// write env file if needed
		var envs []string
		envs, err = writeEnvFilesAndGetEnvs(b)
		if err != nil {
			return 3, err
		}

		// lay and manage error
		exitStatus, err := b.Module.Exec(b, "remove", conf.OtherOptions, envs)
		if err != nil {
			skipFollowing = true
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = 3
		} else if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("remove return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		} else {
			report.Status = TAG_OK
		}

		execSummary[i] = report
	}

	execSummary.Display()
	return
}
