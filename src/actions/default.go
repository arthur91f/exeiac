package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	"strings"

	color "github.com/fatih/color"
)

type ExecSummary []ExecReport

type ExecReport struct {
	Brick *exinfra.Brick
	Error error
}

func (es ExecSummary) Display() {
	var sb strings.Builder

	sb.WriteString(color.New(color.Bold).Sprint("\nSummary:\n"))
	for _, report := range es {
		sb.WriteString("- ")
		if report.Error != nil {
			sb.WriteString(color.New(color.Bold).Sprint("["))
			sb.WriteString(color.RedString("NOK"))
			sb.WriteString(color.New(color.Bold).Sprint("]"))
		} else {
			sb.WriteString(color.New(color.Bold).Sprint("["))
			sb.WriteString(color.GreenString("OK"))
			sb.WriteString(color.New(color.Bold).Sprint("]"))
		}
		sb.WriteString(fmt.Sprintf(" %s", report.Brick.Name))
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
}

func (es ExecSummary) String() string {
	var sb strings.Builder

	sb.WriteString("Summary:\n")
	for _, report := range es {
		if report.Error != nil {
			sb.WriteString("Failed")
		} else {
			sb.WriteString("Succeess")
		}
		sb.WriteString(fmt.Sprintf(" %s", report.Brick.Name))
		sb.WriteString("\n")
	}

	return sb.String()
}

// Triggers a module execution for very single brick in `bricksToExecute`
// Ignores errors and calls the action in `args.Action` for every single brick,
// then prints out a summary of it all.
// Exit code matches 3 if an error occured, 0 otherwise.
func Default(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {
	if infra == nil && args == nil {
		err = exargs.ErrBadArg{Reason: "Error: infra and args are not set"}

		return
	}

	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		_, exitError, err := b.Module.Exec(b, args.Action, args.OtherOptions)

		if exitError != nil {
			statusCode = exitError.ExitCode()
		}

		if err != nil {
			if _, is := err.(exinfra.ActionNotImplementedError); is {
				// NOTE(half-shell): if action if not implemented, we don't take it as an error
				// and move on with the execution
				fmt.Printf("%v ; assume there is nothing to do.\n", err)
				err = nil
			} else {
				statusCode = 3
			}
		}

		execSummary[i] = ExecReport{
			Brick: b,
			Error: err,
		}
	}

	execSummary.Display()

	return
}
