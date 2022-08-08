#!/bin/bash

function get_selected_bricks { #< -brick-path -bricks-paths-list -bricks-names-list
    #> selected_bricks
    brick_path="$(get_arg --string=brick-path "$@")"
    selected_bricks=""
    return_code=0
    if [ -z "$brick_path" ]; then
        case "$(get_brick_type "$brick_path")" in
            not_a_brick)
                soft_exit 1 "ERROR:get_selected_bricks:not_a_brick: $brick_path"
                ;;
            elementary_script_brick|elementary_directory_brick)
                selected_bricks="$bricks_path"
                ;;
            super_brick)
                selected_bricks="$(get_child_bricks "$brick_name")"
                ;;
            *)
                soft_exit 1 "ERROR:get_selected_bricks: brick type not known"
                ;;
        esac
    elif get_arg --string=bricks-paths-list "$@" >/dev/null; then
        selected_bricks="$(get_arg --string=bricks-paths-list "$@" >/dev/null)"
    elif get_arg --string=bricks-names-list "$@" >/dev/null; then
        selected_bricks="$(get_arg --string=bricks-names-list "$@" >/dev/null)"
    fi
    selected_bricks="$(convert_to_elementary_bricks_path "$selected_bricks")"
    if [ "$?" != 0 ]; then
        return_code=1
    fi
    get_bricks_paths_list "$selected_bricks"
    if [ "$?" != 0 ]; then
        return_code=1 
    fi
    return $return_code
}

function get_specified_bricks { #< -selected-bricks -bricks-specifier
    #> specified_bricks_list
    selected_bricks="$(get_arg --string=selected-bricks "$@")"
    specifiers_list="$(get_arg --string=bricks-specifiers "$@")"
    specified_bricks=""
    for specifier in $(sed 's|+|\n|g' <<<"$specifiers_list") ; do
        bricks_to_add=""
        case "$specifier" in
            selected)
                bricks_to_add="$selected_bricks"
                ;;
            all)
                bricks_to_add="$(get_elementary_bricks_list)"
                ;;
            dependencies)
                bricks_to_add="$(get_list_dependencies \
                    "$(get_bricks_paths_list "$selected_bricks")")"
                ;;
            recursive_dependencies)
                bricks_to_add="$(get_list_dependencies_recursively \
                    "$(get_bricks_paths_list "$selected_bricks")")"
                ;;
            dependents)
                bricks_to_add="$(get_list_dependents \
                    "$(get_bricks_names_list "$selected_bricks")")"
                ;;
            recursive_dependents)
                bricks_to_add="$(get_list_dependents_recursively \
                    "$(get_bricks_names_list "$selected_bricks")")"
                ;;
            dependents_dependencies)
                dependencies="$(get_list_dependencies \
                    "$(get_bricks_paths_list "$selected_bricks")")"
                bricks_to_add="$(get_list_dependents \
                    "$(get_bricks_names_list "$dependencies")")"
                ;;
            recursive_dependents_dependencies)
                dependencies="$(get_list_dependencies \
                    "$(get_bricks_paths_list "$selected_bricks")")"
                bricks_to_add="$(get_list_dependents_recursively \
                    "$(get_bricks_names_list "$dependencies")")"
                ;;
            *)
                soft_exit 1 "ERROR:bad_specifier:\"$specifier\""
                ;;
        esac
        if [ "$?" != 0 ]; then
            return_code=1
        fi
        specified_bricks="$(merge_string_on_new_line
            "$specified_bricks" "$bricks_to_add")"
    done
    display_bricks_in_right_order "$(sort -u <<<"$specified_bricks")"
    if [ "$?" != 0 ]; then
        return_code=1 
    fi
    return $return_code
}

