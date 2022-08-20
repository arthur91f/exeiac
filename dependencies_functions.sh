#!/bin/bash

function get_dependencies_tree { #< [brick_path]
    #> bricks_path_ordered_list_with_their_dependencies # "brick1:deps1 deps2", "brick2:deps1 deps3 deps4"
    if [ -n "$1" ]; then
        brick_path="$1"
        bricks_paths_list="$( display_line_after_match \
            "$(get_all_bricks_paths)" "$brick_path")"
    else
        bricks_paths_list="$(get_all_bricks_paths)"
    fi
    
    for brick_path in $bricks_paths_list ; do
        line="$brick_path: "
        echo "$brick_path:$( echo $(get_dependencies "$brick_path"))"
    done
}

function get_dependents { #< brick_path
    #> bricks_paths_ordered_list
    brick_path="$1"
    return_code=0
    bricks_to_check="$( display_line_after_match \
        "$(get_all_bricks_paths)" "$brick_path")"
    
    for brick in $bricks_to_check ; do
        dependencies_list="$(get_dependencies "$brick")"
        if [ "$?" != 0 ]; then
            echo "ERROR:get_dependents:get_dependencies $brick" >&2
            return_code=1
        fi
        if grep -q "^$brick_path$" <<<"$dependencies_list"; then
            echo "$brick"
        fi
    done
    return $return_code
}

function get_dependents_recursively { #< brick_path [-dependencies-tree]
    #> bricks_paths_ordered_list
    brick_path="$1"
    dependents_list="$brick_path"
    if ! dependencies_tree="$(get_arg --string=dependencies-tree "$@")"; then
        dependencies_tree="$(get_dependencies_tree "$brick_path")"
    fi
    while read line ; do
        studied_brick="$(cut -d: -f1 <<<"$line")"
        studied_bricks_dependencies="$(cut -d: -f2- <<<"$line")"
        for dependency in $studied_bricks_dependencies ; do
            if grep -q "^$dependency$" <<<"$dependents_list"; then
                dependents_list="$dependents_list $dependency"
                echo "$studied_brick"
                break
            fi
        done
    done <<<"$dependencies_tree"
}

function get_list_dependents { #< bricks_paths_list
    #> bricks_disordered_list
    bricks_list="$1"
    return_code=0
    dependents_list=""
    for brick in $bricks_list ; do
        get_dependents "$brick"
        if [ "$?" != 0 ]; then
            echo "ERROR:get_list_dependents:$brick"
            return_code=1
        fi
    done | sort -u
    return $return_code
}

function get_list_dependents_recursively { #< bricks_paths_list
    #> bricks_disordered_list
    bricks_list="$(display_bricks_in_right_order "$1")"
    dependencies_tree="$(get_dependencies_tree "$(head -n1 <<<"bricks_list")")"
    return_code=0
    for brick in $brick_list ; do
        get_dependents_recursively "$brick" \
            -dependencies-tree "$dependencies_tree"
        if [ "$?" != 0 ]; then
            return_code=1
        fi
    done | sort -u # sort -u use the right order 
    return $return_code
}

function get_dependencies { #< brick_path [-brick-type]
    #> bricks_paths_disorder_list
    return_code=0
    if brick_type="$(get_arg --string=brick-type)"; then
        brick_names_list="$(execute_brick \
            -action=show_dependencies \
            -brick-path="$1" \
            -brick-type="$brick_type")"
    else
        brick_names_list="$(execute_brick \
            -action=show_dependencies -brick-path="$1")"
    fi
    if [ "$?" != 0 ]; then
        echo "ERROR:get_dependencies:execute_brick show dependencies $1"
        return_code=1
    fi
    get_bricks_paths_list "$brick_names_list"
    if [ "$?" != 0 ]; then
        echo "ERROR:get_dependencies:convert_to_brick_path $1"
        return_code=1
    fi
    return $return_code
}

function get_dependencies_recursively { #< brick_path
    #> bricks_disorder_list
    brick_path="$1"
    dependencies_list="$(get_dependencies "$brick_path")" # fill step by step
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

function get_list_dependencies { #< bricks_paths_list
    #> bricks_disorder_list
    bricks_list="$1"
    return_code=0
    dependencies_list=""
    for brick in $bricks_list ; do
        if brick_type=$(get_brick_type "$brick") ; then
            dependencies_list="$(merge_string_on_new_line \
                "$dependencies_list"\
                "$(get_dependencies "$brick" -brick-type="$brick_type")")"
            if [ "$?" != 0 ]; then
                echo "ERROR:get_dependencies_list:fail on $brick" >&2
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

function get_list_dependencies_recursively { #< bricks_paths_list
    #> bricks_disorder_list
    bricks_paths_list="$1"
    return_code=0
    for brick in $bricks_paths_list ; do
        get_dependencies_recursively "$brick"
        if [ "$?" != 0 ]; then
            echo "ERROR:get_list_dependencies_recursively: $brick" >&2
            return_code=1
        fi
    done | sort -u
    return $return_code
}

