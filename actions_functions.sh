#!/bin/bash

function execute_brick {
    case "$(get_brick_type "$brick_path")" in
        super_brick)
            list_bricks="$(list_bricks "$brick_name")"
            for sub_brick in $(ls -1 "$brick_path") ; do
                if get_brick_type "$brick_path/$sub_brick" >/dev/null ; then
                    echo "## EXEC BRICK: $(get_brick_name "$brick_path/$sub_brick")"
                    $0 "$brick_path/$sub_brick" "$action" "$@"
                    if [ "$?" != 0 ]; then
                        echo "ERROR: $action fail on $(get_brick_name "$brick_path")"
                        exit 1
                    fi
                    echo "## ---------------- ##"
                fi
            done
        ;;
        elementary_script_brick)
            cd "$(dirname "$brick_path")"
            source "$exeiac_lib_path/default_module.sh"
            source "$brick_path"
            $action "$@"
            cd "$current_path"
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
                echo "Error: no known modules can execute that brick: $brick_name" >&2
                exit 5
            fi
            $action "$@"
            cd "$current_path"
        ;;
        not_a_brick)
            echo "Error: not a brick: $brick_name" >&2
            exit 3
        ;;
        *)
            echo "Error: brick with unrecognized type: $brick_name" >&2
            exit 4
        ;;
    esac
}

function list_bricks {
    brick_name="$1"
    if [ -n "$brick_name" ]; then
        elementary_bricks_list="$(get_elementary_bricks_list |
            grep "^$brick_name")"
    else        
        elementary_bricks_list="$(get_elementary_bricks_list)"
    fi              
    echo "$elementary_bricks_list" 
}

function show_dependents {
    bricks_to_check="$(
        display_line_after_match "$(get_elementary_bricks_list)" "$brick_name")"
    for brick in $bricks_to_check ; do
        if grep -q "$brick_name" <<<"$($0 "$brick" show_dependencies)"; then
            echo "$brick"
        fi
    done
}

function display_help {
    cat "$exeiac_lib_path/help.txt"
}

