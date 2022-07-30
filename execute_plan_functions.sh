#!/bin/bash

function get_dependencies_tree_after_brick {
    bricks_list="$(
        display_line_after_match "$(get_elementary_bricks_list)"  "$brick_name")"
    
    for brick in $bricks_list ; do
        echo $brick_name:$($0 "$brick" show_dependencies)
    done
}

function get_dependents_recursively {
    dependents_list="$brick_name"
    get_dependencies_tree_after_brick | while read line ; do
        studied_brick="$(cut -d: -f1 <<<"$line")"
        studied_bricks_dependencies="$(cut -d: -f2- <<<"$line")"
        for dependency in studied_bricks_dependencies ; do
            if grep -q "^$dependency$" <<<"$dependents_list"; then
                dependents_list="$dependents_list $dependency"
                echo "$studied_brick"
                break
            fi
        done
    done
}

# function get_dependencies_tree_before_brick {
#    bricks_list="$(
#        display_line_before_match "$(get_elementary_bricks_list)"  "$brick_name")"
#    for brick in $bricks_list ; do
#        echo $brick_name:$($0 "$brick" show_dependencies)
#    done 
#}

function display_brick_in_right_order {
    brick_to_display="$1"
    get_elementary_bricks_list | while read brick ; do
        grep "^$brick$" <<<"$brick_to_display"
    done
    return 0
}

function get_dependencies_recursively {
    dependencies_list="$($0 "$brick_path" show_dependencies)"
    just_added_dependencies="$dependencies_list"
    while [ -n "$just_added_dependencies" ]; do
        
        # search the dependency of just_added_dependencies
        new_or_already_knowed_dependencies=""
        for dependency in $just_added_dependencies; do
            new_or_already_knowed_dependencies="$(merge_string_on_new_line \
                "$new_or_already_knowed_dependencies" \
                "$($0 "$dependency" show_dependencies)")"
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
            display_brick_in_right_order "$dependencies_list"
        fi
    done
}

function plan_bricks_list {
    bricks_list="$1" # expect they are in the right order
    plan_sum_up=""
    for brick in $bricks_list ; do
        $0 "$dependency" plan
        case "$?" in
            0)
                plan_sum_up="$(
                    merge_string_on_new_line "$plan_sum_up" "ok:$brick")"
            ;;
            1)
                echo "Error:plan: $brick" >&2
                plan_sum_up="$(
                    merge_string_on_new_line "$plan_sum_up" "error1:$brick")"
            ;;
            2)
                plan_sum_up="$(
                    merge_string_on_new_line "$plan_sum_up" "error1:$brick")"
            ;;
            *)
                echo "Error:plan: non expected return code $?:$brick" >&2
                plan_sum_up="$(
                    merge_string_on_new_line "$plan_sum_up" "error$?:$brick")"
            ;;
        esac
    done
    echo "$plan_sum_up"
    if grep -q "^error[0-9]*:" <<<"$plan_sum_up"; then
        return 1
    elif grep -q "^diff:" <<<"$plan_sum_up"; then
        return 2
    else
        return 0
    fi
}

function apply_bricks_list {
    bricks_list="$1"
    apply_sum_up=""
    brick_apply_failed=""
    for brick in $bricks_list ; do
        $0 "$dependency" apply
        if [ "$?" == 0 ]; then
            apply_sum_up="$(
                merge_string_on_new_line "$apply_sum_up" "apply:ok:$brick")"
        else
            apply_sum_up="$(
                merge_string_on_new_line "$apply_sum_up" "apply:error:$brick")"
            brick_apply_failed="$brick"
            break
        fi
    done
    if [ -n "$brick_apply_failed" ]; then
        plan_bricks_list "$(
            display_line_after_match "bricks_list" "$brick_apply_failed")" |
            sed 's|^\(.*\)$|plan:\1|g'
        echo "$apply_sum_up"
        return 1
    else
        echo "$apply_sump_up"
        return 0
    fi
}

function destroy_bricks_list {
    bricks_list="$(tac <<<"$1")"
    destroy_sum_up=""
    brick_destroy_failed=""
    for brick in $bricks_list ; do
        $0 "$dependency" destroy
        if [ "$?" == 0 ]; then
            destroy_sum_up="$(
                merge_string_on_new_line "$destroy_sum_up" "destroy:ok:$brick")"
        else
            destroy_sum_up="$(
                merge_string_on_new_line "$destroy_sum_up" "destroy:error:$brick")"
            brick_destroy_failed="$brick"
            break
        fi
    done
    if [ -n "$brick_destroy_failed" ]; then
        destroys_skipped="$(
            display_line_after_match "bricks_list" "$brick_destroy_failed" | 
            grep -v "$brick_destroy_failed")"
        echo "$destroy_sum_up"
        echo "$destroys_skipped" | sed 's|^\(.*\)$|destroy:skipped:\1|g'
        return 1
    else
        echo "$destroy_sum_up"
        return 0
    fi
}



