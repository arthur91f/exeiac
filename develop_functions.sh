#!/bin/bash
# Here are function for debug and unit testing
function dispdebug { #< string_to_display_on_error_output [tag] #>2
    return_code=$? # to be transparent and don't change return code
    if [ -n "$2" ]; then
        tag="$2"
    else
        tag="debug:"
    fi
    number_line="$(wc -l <<<"$1")"
    case "$number_line" in
    0)
        echo "$tag <nothing to display>" >&2
        ;;
    1)
        echo "$tag $1" >&2
        ;;
    *)
        while read line ; do
            echo "$tag $line" >&2
        done <<<"$1"
        ;;
    esac
    return $return_code
}

function cmd_debug { #< -disbale-debug-classic-output -debug-command
    #> ? ~
    if ! get_arg --boolean=disbale-debug-classic-output "$@" ; then
        echo "### DEBUG ###"
        echo "  action: $action"
        echo "  selected_bricks: \"$selected_bricks\""
        echo "  execute_plan: \"$execute_plan\""
        echo "  brick_path: $brick_path"
        echo "  brick_name: $brick_name"
        echo "  OPTS: ${OPTS[@]}"
        echo "#############"
    fi
    if cmd="$(get_arg --string=debug-command "$@")"; then
        eval "$cmd"
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

