package actions

import (
	"bytes"
	"fmt"
	"os"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Lay(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {

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
			report.Status = "SKIP"
			execSummary[i] = report
			continue
		}

		// write env file if needed
		formatters, envFormatter, err := b.CreateFormatters()
		if err != nil {
			return 3, err
		}

		if len(formatters) > 0 {
			for path, formatter := range formatters {
				f, err := os.Create(path)
				if err != nil {
					return 3, err
				}

				data, err := formatter.Format()
				_, err = f.Write(data)
				if err != nil {
					return 3, err
				}
			}
		}

		// lay and manage error
		exitStatus, err := b.Module.Exec(b, "lay", []string{}, envFormatter.Environ())
		if err != nil {
			skipFollowing = true
			report.Error = err
			report.Status = "ERR"
		} else if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("lay return: %b", exitStatus)
			report.Status = "ERR"
		}

		// check if outputs has changed
		stdout := exinfra.StoreStdout{}
		exitStatus, err = b.Module.Exec(b, "output", []string{}, envFormatter.Environ(), &stdout)
		if err != nil {
			skipFollowing = true
			report.Error = fmt.Errorf("layed apparently success but output failed : %v", err)
			report.Status = "ERR"
		}
		if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("layed apparently success but output return : %b", exitStatus)
			report.Status = "ERR"
		}
		if bytes.Compare(stdout.Output, b.Output) == 0 {
			report.Status = "OK"
		} else {
			b.Output = stdout.Output
			report.Status = "DONE"
		}
		execSummary[i] = report

	}

	execSummary.Display()
	for _, report := range execSummary {
		if !(report.Status == "OK" || report.Status == "DONE") {
			statusCode = 3
		}
	}
	return
}
