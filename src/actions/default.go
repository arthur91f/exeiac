package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
	extools "src/exeiac/tools"
)

// Triggers a module execution for very single brick in `bricksToExecute`
// Ignores errors and calls the action in `args.Action` for every single brick,
// then prints out a summary of it all.
// Exit code matches 3 if an error occured, 0 otherwise.
func PassthroughAction(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		return exstatuscode.INIT_ERROR,
			exinfra.ErrBadArg{Reason: fmt.Sprintf("Error: you should specify at least a brick for %s action", conf.Action)}
	}

	execSummary := make(ExecSummary, len(bricksToExecute))

	var bricksToOutput exinfra.Bricks
	bricksToOutput, err = getBricksToOutput(bricksToExecute, infra, conf.Action)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
	}

	err = enrichOutputs(bricksToOutput)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
	}

	for i, b := range bricksToExecute {

		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		envs, err := writeEnvFilesAndGetEnvs(b, conf.Action)
		if err != nil {
			statusCode = exstatuscode.Update(statusCode, exstatuscode.RUN_ERROR)
			report.Error = fmt.Errorf("not able to get env file and vars before execute: %v", err)
			report.Status = TAG_ERROR
			continue
		}

		events, err := b.Module.Exec(b, conf.Action, conf.OtherOptions, envs)

		if err != nil {
			if actionNotImplementedError, isActionNotImplemented := err.(exinfra.ActionNotImplementedError); isActionNotImplemented {
				// NOTE(half-shell): if action if not implemented, we don't take it as an error
				// and move on with the execution
				fmt.Printf("%v ; assume there is nothing to do.\n", actionNotImplementedError)
				err = nil
				report.Status = "OK"
			} else {
				statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
				report.Status = "ERR"
				report.Error = err
			}
		} else if len(events) > 0 {
			report.Error = fmt.Errorf("WARNING: events aren't yet consumed by exeiac for clean action")
			report.Status = "DONE"
		} else {
			report.Status = "DONE"
		}

		execSummary[i] = report
		fmt.Println("")
	}

	execSummary.Display()

	return
}
