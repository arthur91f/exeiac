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

		// plan and manage error
		events, err := b.Module.Exec(b, "plan", conf.OtherOptions, envs)
		if err != nil {
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
		} else if len(events) == 0 {
			report.Status = TAG_MAY_DRIFT
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT_OR_NOT)
		} else if exeiac_plan, isPresent := events["exeiac_plan"]; isPresent {
			if exeiac_plan == "no_drift" {
				report.Status = TAG_NO_CHANGE
			} else if exeiac_plan == "drift" {
				report.Status = TAG_DRIFT
				statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT)
			} else if exeiac_plan == "unknown" {
				report.Status = TAG_MAY_DRIFT
				statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT_OR_NOT)
			} else {
				report.Error = fmt.Errorf("%s events exeiac_plan has unrecognize value: %s",
					b.Name, exeiac_plan)
				report.Status = TAG_ERROR
				statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_ERROR)
			}
		} else if exeiac_plan_no_drift, isPresent := events["exeiac_plan_no_drift"]; isPresent && exeiac_plan_no_drift == true {
			report.Status = TAG_NO_CHANGE
		} else if exeiac_plan_drift, isPresent := events["exeiac_plan_drift"]; isPresent && exeiac_plan_drift == true {
			report.Status = TAG_DRIFT
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT)
		} else if exeiac_plan_unkown, isPresent := events["exeiac_plan_unkown"]; isPresent && exeiac_plan_unkown == true {
			report.Status = TAG_MAY_DRIFT
			statusCode = exstatuscode.Update(statusCode, exstatuscode.MODULE_DRIFT_OR_NOT)
		}

		execSummary[i] = report
		fmt.Println("")
	}

	execSummary.Display()

	return
}
