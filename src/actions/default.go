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

	for i, b := range bricksToExecute {

		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		// NOTE(arthur91f): we may need to add:    envs, err := writeEnvFilesAndGetEnvs(b)
		// it seems not necessary for init and validate code but who knows for other actions
		exitStatus, err := b.Module.Exec(b, conf.Action, conf.OtherOptions, []string{})

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
		} else if exitStatus == 0 {
			report.Status = "DONE"
		} else {
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
			report.Status = "ERR"
			report.Error = fmt.Errorf("module exit with status code %d", exitStatus)
		}

		execSummary[i] = report
	}

	execSummary.Display()

	return
}
