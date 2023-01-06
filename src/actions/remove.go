package actions

import (
	"fmt"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
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

		return exstatuscode.INIT_ERROR, err
	}

	if conf.Interactive {
		fmt.Println("Here, the bricks list to remove :")
		fmt.Print(bricksToExecute)

		// NOTE(half-shell): We might change this behavior to only ask for a "\n" input
		// instead of a Y/N choice.
		confirm, err := extools.AskConfirmation("\nDo you want to continue ?")

		if err != nil {
			return exstatuscode.RUN_ERROR, err
		} else if !confirm {
			return 0, nil
		}
	}

	err = enrichDatas(bricksToExecute, infra)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
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
			fmt.Printf("remove skipped\n\n")
			continue
		}

		// write env file if needed
		var envs []string
		envs, err = writeEnvFilesAndGetEnvs(b)
		if err != nil {
			statusCode = exstatuscode.Update(statusCode, exstatuscode.RUN_ERROR)
			report.Error = fmt.Errorf("not able to get env file and vars before execute: %v", err)
			report.Status = TAG_ERROR
			continue
		}

		// lay and manage error
		exitStatus, err := b.Module.Exec(b, "remove", conf.OtherOptions, envs)
		if err != nil {
			skipFollowing = true
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
		} else if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("remove return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
		} else {
			report.Status = TAG_OK
		}

		execSummary[i] = report
		fmt.Println("")
	}

	execSummary.Display()
	return
}
