package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
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
	if len(conf.BricksNames) == 0 {
		fmt.Println(defaultHelp)
		return
	} else {
		// TODO(arthur91f): let browse the execute plan to check all different bricks help
		return 3, exinfra.ErrBadArg{
			Reason: "Help action take 0 or one brick as arg, not more"}
	}
	/*
		-- infra-ground/envs/production/network --
		help: no specific help for this module
		--  infra-ground/envs/production/bastion --
		default help with an additionnal action:
		format: rewrite your file: terraform fmt, go fmt -w...
	*/
}
