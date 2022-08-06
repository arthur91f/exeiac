#!/bin/bash

function convert_to_elementary_bricks_path { # bricks_list
    bricks_list="$(get_bricks_paths_list "$1")"
    return_code=0
    for brick in $bricks_list; do
        brick_type="$(get_brick_type "$brick")"
        if [ "$?" != 0 ]; then
            echo "ERROR:convert_to_elementary_bricks_path:get_brick_type:$brick"
            return_code=1
        elif [ "$brick_type" == "super_brick" ]; then
            list_bricks "$(get_brick_name "$brick")"        
        else
            echo "$brick"
        fi
    done
    return $return_code
}

function write_sum_up { # string
    echo "$1" >> "$EXECUTE_SUM_UP_FILE"
}

function get_elementary_bricks_list { # [-rooms-list]
    if ! rooms_list="$(get_arg -string=rooms-list)" ; then
        rooms_list=$ROOMS_LIST
    fi
    bricks_path_list="$(for room_path in $rooms_list ; do
        cd "$room_path"
        find . | grep "/[0-9]\+-[^/]*$" | grep -v '/[^0-9]' |
        sed "s|^\./|$room_path/|g" ; done)"
    for brick_path in $bricks_path_list ; do
        if [ "$(get_brick_type "$brick_path")" != "super_brick" ]; then
            get_brick_name "$brick_path"
        fi
    done
}

function get_dependencies_tree { # [brick_name]
    if [ -n "$1" ]; then
        brick_name="$1"
        bricks_list="$( display_line_after_match \
            "$(get_elementary_bricks_list)" \
            "$brick_name")"
    else
        bricks_list="$(get_elementary_bricks_list)"
    fi
    
    for brick in $bricks_list ; do
        echo "$brick:$( echo $(execute_brick \
            -action=show_dependencies \
            -brick-path=$brick))"
    done
}

function get_dependents_recursively { # brick_name [-dependencies-tree]
    brick_name="$1"
    dependents_list="$brick_name"
    if ! dependencies_tree="$(get_arg --string=dependencies-tree "$@")"; then
        dependencies_tree="$(get_dependencies_tree "$brick_name")"
    fi
    while read line ; do
        studied_brick="$(cut -d: -f1 <<<"$line")"
        studied_bricks_dependencies="$(cut -d: -f2- <<<"$line")"
        for dependency in studied_bricks_dependencies ; do
            if grep -q "^$dependency$" <<<"$dependents_list"; then
                dependents_list="$dependents_list $dependency"
                echo "$studied_brick"
                break
            fi
        done
    done <<<"$dependencies_tree"
}

function display_bricks_in_right_order { # brick_to_display
    brick_to_display="$1"
    get_elementary_bricks_list | while read brick ; do
        grep "^$brick$" <<<"$brick_to_display"
    done
    return 0
}

function get_dependencies_recursively { # brick_path
    brick_path="$1"
    dependencies_list="$(execute_brick \
        -action=show_dependencies \
        -brick-path="$brick_path")"
    just_added_dependencies="$dependencies_list"
    
    while [ -n "$just_added_dependencies" ]; do
        
        # search the dependency of just_added_dependencies
        new_or_already_knowed_dependencies=""
        for dependency in $just_added_dependencies; do
            new_or_already_knowed_dependencies="$(merge_string_on_new_line \
                "$new_or_already_knowed_dependencies" \
                "$(get_dependencies "$dependency")" )"
        done

        # search if dependencies found are not already added
        dependencies_to_add=""
        for dependency in $new_or_already_knowed_dependencies ; do
            if ! grep -q "$dependency" <<<"$dependencies_list" && 
                ! grep -q "$dependency" <<<"$dependencies_to_add" ; then
                dependencies_to_add="$(merge_string_on_new_line \
                    "$dependencies_to_add" "$dependency")"
            fi
        done

        # add new dependencies and
        dependencies_list="$(merge_string_on_new_line \
                "$dependencies_list" "$dependencies_to_add")"
        just_added_dependencies="$dependencies_to_add"
        if [ -z "$dependencies_to_add" ]; then
            display_bricks_in_right_order "$dependencies_list"
        fi

    done
}

function execute_bricks_list { # -bricks-paths-list -action
    bricks_list="$(get_arg --string=bricks-paths-list "$@")"
    action="$(get_arg --string=action "$@")"
    return_code=0
    case "$action" in
        init|validate|fmt|pass|help) # run all even if it fails
            for brick in $bricks_list ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    echo "## EXEC BRICK: $brick_name"
                    execute_brick -brick-path="$brick" -action=output \
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
        output) # run all even if it fails without exec_sum_up
            echo "{"
            for brick in $bricks_list ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    echo "  \"$brick_name\": {"
                    execute_brick -brick-path="$brick" -action=output \
                        -brick-type="$brick_type" "$@" | sed 's|^\(.*\)$|    \1|g'
                    if [ "$?" !=0 ]; then
                        return_code=1
                        echo "ERROR:$brick_name:output" >&2
                        echo "    \"error\": \"don't manage to output brick\""
                    fi
                    echo "  }"
                else
                    return_code=1
                fi
            done
            echo "}"
        ;;
        show_dependencies) # run all even if it fails without exec_sum_up
            dependencies_list="$(get_list_dependencies "$bricks_list")"
            return_code="$?"
            display_bricks_in_right_order "$dependencies_list"
            if [ "$?" !=0 ]
                return_code=1
            fi
        ;;
        plan)
            for brick in $bricks_list ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    echo "## EXEC BRICK: $brick_name"
                    execute_brick -brick-path="$brick" -action="$action" \
                        -brick-type="$brick_type" "$@"
                    if [ "$?" == 0 ]; then
                        write_sum_up "plan:OK:$brick_name"
                    elif [ "$?" == 1 ]; then
                        echo "ERROR:execute_bricks_list:plan:$brick_name" >&2
                        write_sum_up "plan:ERROR:$brick_name"
                        return_code=1
                    elif [ "$?" == 2 ]; then
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
        apply) # run all until it fails
            for brick in $bricks_list ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    if [ "$return_code" == 0 ]; then
                        echo "## EXEC BRICK: $brick_name"
                        execute_brick -brick-path="$brick" -action="$action" \
                            -brick-type="$brick_type" "$@"
                        if [ "$?" != 0 ]; then
                            echo "ERROR: $action fail on $(get_brick_name "$arg_brick")" >&2
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
        destroy) # run all in reverse until it fails
            for brick in $(tac <<<"$bricks_list") ; do
                brick_name="$(get_brick_name "$brick")"
                if brick_type=$(get_brick_type "$brick") ; then
                    if [ "$return_code" == 0 ]; then
                        echo "## EXEC BRICK: $brick_name"
                        execute_brick -brick-path="$brick" -action="$action" \
                            -brick-type="$brick_type" "$@"
                        if [ "$?" != 0 ]; then
                            echo "ERROR: $action fail on $brick_name" >&2
                            write_sum_up "$action:ERROR:$brick_name"
                        fi
                        echo "## ---------------- ##"
                    else
                        write_sum_up "$action:CANCEL:$brick_name"
                    fi
                fi
            done
        ;;
        *)
            echo "ERROR:execute_bricks_list:unrecognized action: $action" >&2
            return_code=1
        ;;
    esac
    return $return_code
}

function get_dependencies { # brick_path
    execute_brick -action=show_dependencies -brick-path="$1"
    if [ "$?" != 0 ]; then
        echo "ERROR:get_dependencies: $1"
        return 1
    fi
}

function get_list_dependencies { # bricks_paths_list
    bricks_list="$1"
    return_code=0
    dependencies_list=""
    for brick in $bricks_list ; do
        if brick_type=$(get_brick_type "$brick") ; then
            dependencies_list="$(merge_string_on_new_line \
                "$dependencies_list" \ 
                "$(execute_brick -brick-path="$brick" \
                    -action=show_dependencies \
                    -brick-type="$brick_type")")"
            if [ "$?" != 0 ]; then
                echo "ERROR:get_dependencies_list:fail on $(get_brick_name "$brick")" >&2
                return_code=1
            fi
        else
            echo "ERROR:get_dependencies_list:$brick" >&2
                return_code=1
        fi
    done
    sort -u <<<"$dependencies_list"
    return $return_code
}

function get_list_dependencies_recursively { # bricks_paths_list
    bricks_paths_list="$1"
    return_code=0
    for brick in $bricks_paths_list ; do
        get_dependencies_recursively "$brick"
        if [ "$?" != 0 ]; then
            echo "ERROR:get_list_dependencies_recursively: $brick" >&2
            return_code=1
        fi
    done
    return $return_code
}

function get_dependents { # brick_name
    brick_name="$1"
    return_code=0
    bricks_to_check="$( display_line_after_match \
        "$(get_elementary_bricks_list)" "$brick_name")"
    
    for brick in $bricks_to_check ; do
        dependencies_list="$(get_dependencies "$brick")"
        if [ "$?" != 0 ]; then
            echo "ERROR:get_dependents: $brick" >&2
            return_code=1
        fi
        if grep -q "$brick_name" <<<"$dependencies_list"; then
            echo "$brick"
        fi
    done
    return $return_code
}

function get_list_dependents { # bricks_names_list
    bricks_list="$1"
    return_code=0
    dependents_list=""
    for brick in $bricks_list ; do
        get_dependents "$bricks"
        if [ "$?" != 0 ]; then
            echo "ERROR:get_list_dependents:$brick"
            return_code=1
        fi
    done | sort -u
    return $return_code
}

function get_list_dependents_recursively { # bricks_names_list
    bricks_list="$(display_bricks_in_right_order "$1")"
    dependencies_tree="$(get_depencies_tree "$(head -n1 <<<"bricks_list")")"
    return_code=0
    for brick in $brick_list ; do
        get_dependents_recursively "$brick" \
            -dependencies-tree "$dependencies_tree"
        if [ "$?" != 0 ]; then
            return_code=1
        fi
    done | sort -u
    return $return_code
}

