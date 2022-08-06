#!/bin/bash

function list_bricks { # [brick_name]
    brick_name="$1"
    if [ -n "$brick_name" ]; then
        elementary_bricks_list="$(get_elementary_bricks_list |
            grep "^$brick_name")"
    else        
        elementary_bricks_list="$(get_elementary_bricks_list)"
    fi              
    echo "$elementary_bricks_list" 
}

function execute_brick { # -brick-path -action [-brick-type]
    arg_brick="$(get_arg --string=brick-path "$@")"
    arg_action="$(get_arg --string=action "$@")"
    if ! brick_type="$(get_arg --string=brick-type)" ; then
        brick_type="$(get_brick_type "$arg_brick")"        
    fi
    case "$brick_type" in
        super_brick)
            list_bricks="$(list_bricks "$brick_name")"
            execute_bricks_list \
                -bricks-paths-list="$list_bricks" \
                -action="$arg_action"
            return $?
        ;;
        elementary_script_brick)
            cd "$(dirname "$arg_brick")"
            source "$exeiac_lib_path/default_module.sh"
            source "$arg_brick"
            $action "$@"
            return $?
        ;;
        elementary_directory_brick)
            cd "$arg_brick"
            module_found="false"
            for module in $( ls -1 "$modules_path") ; do
                source "$exeiac_lib_path/default_module.sh"
                source "$modules_path/$module"
                if is_brick_using_this_module ; then
                    module_found="true"
                    break
                fi
            done
            if [ "$module_found" == "false" ]; then
                echo "Error: no known modules can execute that brick: $brick_name" >&2
                return 5
            fi
            $action "$@"
            return $?
        ;;
        not_a_brick)
            echo "Error: not a brick: $brick_name" >&2
            return 3
        ;;
        *)
            echo "Error: brick with unrecognized type: $brick_name" >&2
            return 4
        ;;
    esac
}

function show_dependents { #-brick-name
    if [ -z "$1" ]; then
        arg_brick="$(get_arg -string=brick-name)"
    else
        arg_brick="$BRICK_NAME"
    fi
    
    get_dependents "$arg_brick"
    return $?
}

function display_help {
    cat "$exeiac_lib_path/help.txt"
}

