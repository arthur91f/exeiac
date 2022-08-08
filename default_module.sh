#!/bin/bash
function is_brick_using_this_module { #> ?
    echo "default-module: get_type function have to be overloaded" >&2
    return 1
}

function install { #> ~
    echo "default-module: nothing to install"
}

function init { #> ~
    echo "default-module: no init needed"
}

function plan { #> ~
    echo "default-module: no plan possible"
}

function apply { #> ~
    echo "default-module: apply function have to be overloaded" >&2
    return 1
}

function output { #> ~
    echo ""
}

function destroy { #> ~
    echo "default-module: destroy function have to be overloaded" >&2
    return 1
}

function validate { #> ?
    return 0
}

function fmt { #> ?
    return 0
}

function show_dependencies { #> dependencies_list
    egrep -r '(#|//)EXEIAC:depends:' |
        sed 's|^.*#EXEIAC:depends: *||g' |
        sed 's|^.*//EXEIAC:depends: *||g' |
        cut -d: -f1
}

function help { #> ~
    echo "The help of this brick haven't been overloaded 
    so it runs the normal way."
    echo execiac BRICK_PATH ACTION [OPTIONS]
    echo "init: get some dependencies, typically download
    terraform modules or ansible deps"
    echo "plan: a dry run to check what we want to apply"
    echo "apply: run the IaC with the right tools"
    echo "output: display some outputs that can be used by
    other bricks"
    echo "destroy: revert the apply"
    echo "validate: validate that the syntaxe is ok"
    echo "fmt: rewrite file to pass linter"
    echo "show_dependencies: get brick_name whom output is
    needed"
    echo 'install: install tools for run the exeiac
    modules. "init" is specific to brick dependencies,
    "install" is specific to the exeiac module'
}

function pass { #> return true
    # used to execute option like plan-before without plan the actual brick
    return 0
}

### Those function can eventually be used in other modules
function copy_function { #< source_function_name new_function_name
    #> # nothing only duplicates function to overload function without loosing them
    source_function_name="$1"
    new_function_name="$2"
    local local_func
    local_func="$(declare -f "$source_function_name")" &&
    eval "function $new_function_name ${local_func#*"()"}"
}

function import_module_functions { #< module_path
    #> for terrform module will import functions as
    #> init -> terraform_init
    #> plan -> terraform_plan
    #> ...
    #> then it will reimport default function
    # https://webdevdesigner.com/q/how-do-i-rename-a-bash-function-587410/
    module="$1"
    shift
    if [ -z "$@" ]; then
        action_list="init plan apply output destroy validate fmt show_dependencies install"
    else
        action_list="$@"
    fi
    source "$modules_path/$module"
    for action in $actions_list ; do
        local local_func
        local_func="$(declare -f "$action")" &&
        eval "function ${module}_$action ${local_func#*"()"}"
    done
    source "$modules_path/default.sh"
}

