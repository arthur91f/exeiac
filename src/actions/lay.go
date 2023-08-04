package actions

import (
	"bytes"
	"fmt"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
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

		return exstatuscode.INIT_ERROR, err
	}

	if conf.Interactive {
		fmt.Println("Here, the bricks list to lay :")
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

	var bricksToOutput exinfra.Bricks
	bricksToOutput, err = getBricksToOutput(bricksToExecute, infra, conf.Action)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
	}

	bricksToOutput = append(bricksToOutput, bricksToExecute...) // to know if there is a diff
	sort.Sort(bricksToOutput)

	err = enrichOutputs(bricksToOutput)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
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
			fmt.Printf("lay skipped\n\n")
			continue
		}

		// write env file if needed
		envs, err := writeEnvFilesAndGetEnvs(b, conf.Action)
		if err != nil {
			statusCode = exstatuscode.Update(statusCode, exstatuscode.RUN_ERROR)
			report.Error = fmt.Errorf("not able to get env file and vars before execute: %v", err)
			report.Status = TAG_ERROR
			skipFollowing = true
			continue
		}

		_, layExitStatus, layErr := b.Module.Exec(b, "lay", conf.OtherOptions, envs)
		stdout := exinfra.StoreStdout{}
		_, outputExitStatus, outputErr := b.Module.Exec(b, "output", []string{}, envs, &stdout)

		// set skipFollowing, report.Status, report.Error and update b.Ouput
		if layErr == nil && layExitStatus == 0 && outputErr == nil && outputExitStatus == 0 { // everything runs well
			if bytes.Compare(stdout.Output, b.Output) == 0 {
				report.Status = TAG_NO_CHANGE
			} else {
				b.Output = stdout.Output
				report.Status = TAG_DONE
			}
		} else { // there is at least one error
			skipFollowing = true
			report.Status = TAG_ERROR
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)

			// simplify the next condition tree
			if layExitStatus != 0 && layErr == nil {
				layErr = fmt.Errorf("exit with status %d", layExitStatus)
			}
			if outputExitStatus != 0 && outputErr == nil {
				outputErr = fmt.Errorf("exit with status %d", outputExitStatus)
			}

			if layErr != nil && outputErr != nil { // 2 errors
				report.Error = fmt.Errorf("2 errors lay and output error: "+
					"{\"lay\": \"%v\", \"output\": \"%v\"}", layErr, outputErr)
			} else if layErr != nil && outputExitStatus == 0 { // 1 error: check if output changed
				if bytes.Compare(stdout.Output, b.Output) == 0 {
					report.Error = fmt.Errorf(
						"lay has failed, output doesn't seem to has changed: %v",
						layErr)
				} else {
					report.Error = fmt.Errorf(
						"lay has failed, output has changed: %v",
						layErr)
					b.Output = stdout.Output
				}
			} else if layExitStatus == 0 && outputErr != nil { // 1 error: can't get output
				report.Error = fmt.Errorf(
					"lay seems to success but the following output failed: %v",
					outputErr)
			}
		}

		execSummary[i] = report
		fmt.Println("")
	}

	execSummary.Display()

	return
}
