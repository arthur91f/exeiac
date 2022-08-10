#!/bin/bash

function list_bricks { #< [brick_path|bricks_list]
    #> bricks_list
    arg="$1"
    if [ -z "$arg" ]; then
        get_elementary_bricks_list
        return $?
    else
        bricks_list="$(convert_to_elementary_bricks_path "$arg")"
        return_code=$?
        get_bricks_names_list "$bricks_list"
        return $return_code
    fi
}

function execute_brick { #< -brick-path -action [-brick-type]
    #> ~ # depending of -action 
    brick_path="$(get_arg --string=brick-path "$@")"
    action="$(get_arg --string=action "$@")"
    if ! brick_type="$(get_arg --string=brick-type "$@")" ; then
        brick_type="$(get_brick_type "$brick_path")"        
    fi
    case "$brick_type" in
        super_brick)
            bricks_list="$(get_child_bricks "$brick_path")"
            execute_bricks_list \
                -bricks-paths-list="$bricks_list" \
                -action="$action"
            return $?
        ;;
        elementary_script_brick)
            cd "$(dirname "$brick_path")"
            source "$exeiac_lib_path/default_module.sh"
            source "$brick_path"
            $action "$@"
            return $?
        ;;
        elementary_directory_brick)
            cd "$brick_path"
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
                echo "ERROR:execute_brick: no known modules can execute that brick: $brick_path" >&2
                return 5
            fi
            $action "$@"
            return $?
        ;;
        not_a_brick)
            echo "ERROR:execute_brick: not a brick: $brick_path" >&2
            return 3
        ;;
        *)
            echo "ERROR:execute_brick: brick with unrecognized type: $brick_name" >&2
            return 4
        ;;
    esac
}

function display_help { #< # nothing
    #> ~ # display the help
    cat "$exeiac_lib_path/help.txt"
}

