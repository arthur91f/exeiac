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
    if ! rooms_list="$(get_arg -string=rooms-list "$@")" ; then
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

function display_bricks_in_right_order { # brick_to_display
    brick_to_display="$1"
    get_elementary_bricks_list | while read brick ; do
        grep "^$brick$" <<<"$brick_to_display"
    done
    return 0
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
            if [ "$?" !=0 ]; then
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

