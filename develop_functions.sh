#!/bin/bash
# Here are function for debug and unit testing
function dispdebug { #< string_to_display_on_error_output #>2
    echo "debug: $1" >&2
}

function cmd_debug { #< -disbale-debug-classic-output -unit-testing-function
    #> ? ~
    if ! get_arg --boolean=disbale-debug-classic-output "$@" ; then
        echo "### DEBUG ###"
        echo "  arg_case: \"$arg_case\", action: \"$action\""
        echo "  brick_path: $brick_path"
        echo "  brick_name: $brick_name"
        echo "  OPTS: ${OPTS[@]}"
        echo "#############"
    fi
    if cmd="$(get_arg --string=unit-testing-function "$@")"; then
        $cmd
        echo "return code: $?"
    fi
}

function get_functions { #< =(-file|-function)
    #> function_arguments
    
    if fx="$(get_arg --string=function "$@")"; then
        files="$(grep -r "^function $fx" | cut -d: -f1)"
    elif files="$(get_arg --string=file "$@")" ; then
        true
    else
        files="$(ls -1 $exeiac_lib_path/ | grep "_functions\.sh$")"
    fi

    for file in $files ; do
        echo "-- $file ------"
        egrep '(^function|#<|#>)' "$exeiac_lib_path/$file"
        echo ""
    done
}

# function install {
#     return_code=0
#     if [ -z "$room_paths_list" ]; then
#         echo "add room_paths_list variable in $HOME/.exeiac"
#         return_code=1
#     fi
#     if [ "$(type -t sanitize_function)" != "function" ]; then
#         echo "add a sanitize_function in $HOME/.exeiac"
#         return_code=1
#     fi
#     if [ -z "$modules_path" ]; then
#         echo "add a modules_path variable in $HOME/.exeiac"
#         return_code=1
#     fi
#     exit $return_code
# }

