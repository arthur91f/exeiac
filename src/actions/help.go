package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

var defaultHelp = `exeiac ACTIONS (BRICK_PATH|BRICK_NAME)[OPTIONS]
ACTIONS:
init: get some dependencies, typically download terraform modules
    or ansible deps
plan: a dry run to check what we want to lay
lay: lay the brick on the wall. Run the IaC with the right tools
remove: remove a brick from your wall to destroy it properly.
validate_code: validate if the syntaxe is ok
help: display this help or the specified help for the brick
cd: change directory but you can use brick name althought path
show: display brick attributes (depends of the format option choosen)
clean: remove all files created by exeiac
OPTIONS:
-I --non-interactive: run without interaction (use especially for ignore
                      confirmation after lay or remove)
-s --bricks-specifier: (selected|previous|following|children|this|
                        recursive-following|recursive-precedents)
-f --format: (name|path|input|output) use with show
`

func Help(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for help action" +
			"\nif you want to display exeiac help use --help option"}
		return 3, err
	}

	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		extools.DisplaySeparator(b.Name + "(" + b.Module.Name + ")")
		report := ExecReport{Brick: b}
		statusCode, err = b.Module.Exec(b, conf.Action, conf.OtherOptions, []string{})

		if err != nil {
			if _, isActionNotImplemented := err.(exinfra.ActionNotImplementedError); isActionNotImplemented {
				// NOTE(half-shell): if action if not implemented, we don't take it as an error
				// and move on with the execution
				fmt.Printf("help: no specific help for this module: %s\n", b.Module.Name)
				originalActions := extools.StrSliceXor([]string{"show_implemented_actions", "init", "plan", "lay", "remove", "validate_code", "clean", "output"}, b.Module.Actions)
				if len(originalActions) != 0 {
					fmt.Printf("Nethertheless some other actions are implemented for this module:%s\n",
						extools.StringListOfString(originalActions))
				}
				err = nil
				report.Status = "OK"
			} else {
				statusCode = 3
				report.Status = "ERR"
				report.Error = err
			}
		} else {
			report.Status = "DONE"
		}
		execSummary[i] = report
	}

	execSummary.Display()
	return
}
