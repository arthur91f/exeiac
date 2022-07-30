#!/bin/bash
# Here are function for debug and unit testing
function dispdebug {
    echo "debug: $1" >&2
}

function cmd_debug {
    if ! get_arg_in_string "disbale-debug-classic-output" "$EXEIAC_OPTS" >/dev/null ; then
        echo "### DEBUG ###"
        echo "  arg_case: \"$arg_case\", action: \"$action\""
        echo "  brick_path: $brick_path"
        echo "  brick_name: $brick_name"
        echo "  EXEIAC_OPTS: $EXEIAC_OPTS"
        echo "  MODULE_OPTS: $MODULE_OPTS"
        echo "#############"
    fi
    if cmd="$(get_arg_in_string "unit-testing-function" "$EXEIAC_OPTS")"; then
        $cmd
        echo "return code: $?"
    fi
}

function install {
    return_code=0
    if [ -z "$room_paths_list" ]; then
        echo "add room_paths_list variable in $HOME/.exeiac"
        return_code=1
    fi
    if [ "$(type -t sanitize_function)" != "function" ]; then
        echo "add a sanitize_function in $HOME/.exeiac"
        return_code=1
    fi
    if [ -z "$modules_path" ]; then
        echo "add a modules_path variable in $HOME/.exeiac"
        return_code=1
    fi
    exit $return_code
}

