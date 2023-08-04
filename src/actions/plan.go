package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	exstatuscode "src/exeiac/statuscode"
	extools "src/exeiac/tools"
)

func Plan(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for plan action"}

		return exstatuscode.INIT_ERROR, err
	}

	var bricksToOutput exinfra.Bricks
	bricksToOutput, err = getBricksToOutput(bricksToExecute, infra, conf.Action)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
	}

	err = enrichOutputs(bricksToOutput)
	if err != nil {
		return exstatuscode.ENRICH_ERROR, err
	}

	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		// write env file if needed
		envs, err := writeEnvFilesAndGetEnvs(b, conf.Action)
		if err != nil {
			statusCode = exstatuscode.Update(statusCode, exstatuscode.RUN_ERROR)
			report.Error = fmt.Errorf("not able to get env file and vars before execute: %v", err)
			report.Status = TAG_ERROR
			continue
		}

		// lay and manage error
		_, exitStatus, err := b.Module.Exec(b, "plan", conf.OtherOptions, envs)
		if err != nil {
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
		} else if exitStatus == 0 {
			report.Status = TAG_NO_CHANGE
		} else if exitStatus == 2 {
			report.Status = TAG_DRIFT
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT)
		} else if exitStatus == 3 {
			report.Status = TAG_MAY_DRIFT
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT_OR_NOT)
		} else {
			report.Error = fmt.Errorf("plan return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
		}

		execSummary[i] = report
		fmt.Println("")
	}

	execSummary.Display()

	return
}
