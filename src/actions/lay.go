package actions

import (
	"bytes"
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Lay(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for lay action"}
		return 3, err
	} else if len(bricksToExecute) > 1 && conf.Interactive {
		fmt.Print("Here, the bricks list to lay :")
		fmt.Print(bricksToExecute)
		var confirm bool
		confirm, err = extools.AskConfirmation("\nDo you want to continue ?")
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

	for i, b := range bricksToExecute {

		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		// skip if an error was encounter before
		if skipFollowing {
			report.Status = TAG_SKIP
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
		exitStatus, err := b.Module.Exec(b, "lay", conf.OtherOptions, envs)
		if err != nil {
			skipFollowing = true
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = 3
		} else if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("lay return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		}

		// check if outputs has changed
		stdout := exinfra.StoreStdout{}
		exitStatus, err = b.Module.Exec(b, "output", []string{}, envs, &stdout)
		if err != nil {
			skipFollowing = true
			report.Error = fmt.Errorf("layed apparently success but output failed : %v", err)
			report.Status = TAG_ERROR
			statusCode = 3
		}
		if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("layed apparently success but output return : %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		}
		if bytes.Compare(stdout.Output, b.Output) == 0 {
			report.Status = TAG_NO_CHANGE
		} else {
			b.Output = stdout.Output
			report.Status = TAG_DONE
		}

		execSummary[i] = report
	}

	execSummary.Display()
	return
}
