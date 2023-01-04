package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Clean(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		return 3, exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for clean action"}
	}

	if conf.Interactive {
		fmt.Println("Here, the bricks list to clean :")
		fmt.Print(bricksToExecute)

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

	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		envs, err := writeEnvFilesAndGetEnvs(b)
		if err != nil {
			return 3, err
		}

		// clean and manage error
		exitStatus, err := b.Module.Exec(b, "clean", conf.OtherOptions, envs)
		if err != nil {
			if actionNotImplementedError, isActionNotImplemented := err.(exinfra.ActionNotImplementedError); isActionNotImplemented {
				// NOTE(half-shell): if action if not implemented, we don't take it as an error
				// and move on with the execution
				fmt.Printf("%v ; assume there is nothing to do.\n", actionNotImplementedError)
				err = nil
			} else {
				report.Error = err
				report.Status = TAG_ERROR
				statusCode = 3
			}
		} else if exitStatus != 0 {
			report.Error = fmt.Errorf("clean return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		}

		err = cleanEnvFiles(b)
		if err != nil {
			if report.Error != nil {
				report.Error = fmt.Errorf("2 errors:\n"+
					"module error: %v\nclean input file error:%v",
					report.Error, err)
				report.Status = TAG_ERROR
				statusCode = 3
			} else {
				report.Error = err
				report.Status = TAG_ERROR
				statusCode = 3
			}
		} else {
			if report.Error != nil {
				report.Error = fmt.Errorf("module error:%v", report.Error)
				report.Status = TAG_ERROR
				statusCode = 3
			} else {
				report.Status = TAG_DONE
			}
		}

		execSummary[i] = report
	}

	execSummary.Display()

	return
}
