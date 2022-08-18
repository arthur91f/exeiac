#!/bin/bash
# I this file you will found single time called functions
# they are all called by exeiac.sh script

function get_selected_bricks { #< -brick-path -bricks-paths-list -bricks-names-list
    #> selected_bricks
    brick_path="$(get_arg --string=brick-path "$@")"
    selected_bricks=""
    return_code=0
    if [ -n "$brick_path" ]; then
        case "$(get_brick_type "$brick_path")" in
            not_a_brick)
                soft_exit 1 "ERROR:get_selected_bricks:not_a_brick: $brick_path"
                ;;
            elementary_script_brick|elementary_directory_brick)
                selected_bricks="$brick_path"
                ;;
            super_brick)
                selected_bricks="$(get_child_bricks "$brick_path")"
                ;;
            *)
                soft_exit 1 "ERROR:get_selected_bricks: brick type not known"
                ;;
        esac
    elif get_arg --string=bricks-paths-list "$@" >/dev/null; then
        selected_bricks="$(get_absolute_paths_list \
            "$(get_arg --string=bricks-paths-list "$@")")"
    elif get_arg --string=bricks-names-list "$@" >/dev/null; then
        selected_bricks="$(get_bricks_paths_list \
            "$(get_arg --string=bricks-names-list "$@")")"
    fi

    convert_to_elementary_bricks_path "$selected_bricks"
    if [ "$?" != 0 ]; then
        return_code=1
    fi
    return $return_code
}

function get_specified_bricks { #< -selected-bricks -bricks-specifiers
    #> specified_bricks_list
    selected_bricks="$(get_arg --string=selected-bricks "$@")"
    specifiers_list="$(get_arg --string=bricks-specifiers "$@")"
    specified_bricks=""
    for specifier in $(sed 's|+| |g' <<<"$specifiers_list") ; do
        bricks_to_add=""
        case "$specifier" in
            selected)
                bricks_to_add="$selected_bricks"
                ;;
            all)
                bricks_to_add="$(get_all_bricks_paths)"
                ;;
            dependencies)
                bricks_to_add="$(get_list_dependencies "$selected_bricks")"
                ;;
            recursive_dependencies)
                bricks_to_add="$(get_list_dependencies_recursively \
                    "$selected_bricks")"
                ;;
            dependents)
                bricks_to_add="$(get_list_dependents "$selected_bricks")"
                ;;
            recursive_dependents)
                bricks_to_add="$(get_list_dependents_recursively \
                    "$selected_bricks")"
                ;;
            dependents_dependencies)
                dependencies="$(get_list_dependencies "$selected_bricks")"
                bricks_to_add="$(get_list_dependents "$dependencies")"
                ;;
            recursive_dependents_dependencies)
                dependencies="$(get_list_dependencies "$selected_bricks")"
                bricks_to_add="$(get_list_dependents_recursively \
                    "$dependencies")"
                ;;
            *)
                soft_exit 1 "ERROR:bad_specifier:\"$specifier\""
                ;;
        esac
        if [ "$?" != 0 ]; then
            return_code=1
        fi
        specified_bricks="$(merge_string_on_new_line \
            "$specified_bricks" "$bricks_to_add")"
    done
    display_bricks_in_right_order "$(sort -u <<<"$specified_bricks")"
    if [ "$?" != 0 ]; then
        return_code=1 
    fi
    return $return_code
}

function execute_bricks_list { #< -action -bricks-list
    #> ? ~
    action="$(get_arg --string=action "$@")"
    bricks_list="$(get_arg --string=bricks-list "$@")"
    return_code=0

    function write_sum_up { #< string
        #>EXECUTE_SUM_UP_FILE 
        echo "$1" >> "$EXECUTE_SUM_UP_FILE"
    }

    case "$action" in
    init|output|validate|fmt|pass|help) # run all even if it fails
            for brick in $bricks_list ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    echo "## EXEC BRICK: $brick_name"
                    execute_brick -brick-path="$brick" -action="$action" \
                        -brick-type="$brick_type" "$@"
                    if [ "$?" == 0 ]; then
                        write_sum_up "$action:OK:$brick_name"
                    else
                        echo "ERROR:execute_bricks_list:$action:$brick_name" >&2
                        write_sum_up "$action:ERROR:$brick_name"
                        return_code=1
                    fi
                    echo "## ---------------- ##"
                fi
            done
        ;;
    plan)
            for brick in $bricks_list ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    echo "## EXEC BRICK: $brick_name"
                    execute_brick -brick-path="$brick" -action="$action" \
                        -brick-type="$brick_type" "$@"
                    return_code=$? # if not set here the if test will change $?
                    if [ "$return_code" == 0 ]; then
                        write_sum_up "plan:OK:$brick_name"
                    elif [ "$return_code" == 1 ]; then
                        echo "ERROR:execute_bricks_list:plan:$brick_name" >&2
                        write_sum_up "plan:ERROR:$brick_name"
                        return_code=1
                    elif [ "$return_code" == 2 ]; then
                        write_sum_up "plan:DRIFT:$brick_name"
                        if [ "$return_code" == 0 ]; then # to keep worst return_code
                            return_code=2
                        fi
                    else
                        echo "ERROR:execute_bricks_list:plan:$brick_name" >&2
                        write_sum_up "plan:ERROR:$brick_name"
                        return_code=1
                    fi
                else
                    echo "ERROR:execute_bricks_list:unrecognize_brick_type:$brick_name" >&2
                    write_sum_up "plan:ERROR:$brick_name"
                    return_code=1
                fi
                echo "## ---------------- ##"
            done
        ;;
    apply|destroy)
        if [ "$action" == 'destroy' ]; then
            bricks_list="$(sort -r <<<"$destroy")"
        fi
        for brick in $bricks_list ; do
            brick_name="$(get_brick_name "$brick")"
            if brick_type=$(get_brick_type "$brick") ; then
                if [ "$return_code" == 0 ]; then
                    echo "## EXEC BRICK: $brick_name"
                    execute_brick -brick-path="$brick" -action="$action" \
                        -brick-type="$brick_type" "$@"
                    if [ "$?" != 0 ]; then
                        echo "ERROR: $action fail on $brick_name" >&2
                        write_sum_up "$action:ERROR:$brick_name"
                        break
                    fi
                    echo "## ---------------- ##"
                else
                    write_sum_up "$action:CANCEL:$brick_name"
                fi
            fi
        done
        ;;
    show_dependencies)
        bricks_list="$(get_list_dependencies "$bricks_list")"
        return_code="$?"
        get_bricks_names_list "$(display_bricks_in_right_order "$bricks_list")"
        if [ "$?" != 0 ]; then
            return_code=1
        fi
        ;;
    show_dependents)
        bricks_list="$(get_list_dependents "$bricks_list")"
        return_code="$?"
        get_bricks_names_list "$(display_bricks_in_right_order "$bricks_list")"
        if [ "$?" != 0 ]; then
            return_code=1
        fi
        ;;
    show_dependencies_recursively)
        bricks_list="$(get_list_dependencies_recursively "$bricks_list")"
        return_code="$?"
        get_bricks_names_list "$(display_bricks_in_right_order "$bricks_list")"
        if [ "$?" != 0 ]; then
            return_code=1
        fi
        ;;
    show_dependents_recursively)
        bricks_list="$(get_list_dependents_recursively "$bricks_list")"
        return_code="$?"
        get_bricks_names_list "$(display_bricks_in_right_order "$bricks_list")"
        if [ "$?" != 0 ]; then
            return_code=1
        fi
        ;;
    list_bricks)
        bricks_list="$(convert_to_elementary_bricks_path "$bricks_list")"
        return_code=$?
        get_bricks_names_list "$bricks_list"
        ;;
    debug)
        cmd_debug "$@"
        return_code=$?
        ;;
    *)
        echo "ERROR:execute_bricks_list:unrecognized_action:$action" >&2
        return_code=1
        ;;
    esac
    return $return_code
}

function display_help { #< # nothing
    #> ~ # display the help
    cat "$exeiac_lib_path/help.txt"
}

