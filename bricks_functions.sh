#!/bin/bash

function execute_brick { #< -brick-path -action [-brick-type]
    #> ~ # depending of -action 
    brick_path="$(get_arg --string=brick-path "$@")"
    action="$(get_arg --string=action "$@")"
    if ! brick_type="$(get_arg --string=brick-type "$@")" ; then
        brick_type="$(get_brick_type "$brick_path")"
    fi
    case "$brick_type" in
        super_brick)
            echo "ERROR:execute_brick:bad_brick_type:super_brick:$brick_path"
            return 1
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

function get_brick_sanitized_name { #< brick_name
    #> sanitize_brick_name
    brick_name="$1"
    sanitize_function "$brick_name"
}

function get_brick_type { #< brick_path
    #> ? if recognize a brick_type
    #> (super_brick|elementary_script_brick|elementary_directory_brick|not_a_brick)
    brick_path="$1"
    if ! grep -q "/[0-9]\+-[^/]*$" <<<"$brick_path" ; then
        echo "not_a_brick"
        return 1
    fi
    if [ -f "$brick_path" ] && [ -x "$brick_path" ] ; then
        echo "elementary_script_brick"
    elif [ -d "$brick_path" ] &&
        ls -1 "$brick_path" | grep -q "^[0-9]\+-[^/]*$" ; then
        echo "super_brick"
    elif [ -d "$brick_path" ]; then
        echo "elementary_directory_brick"
    else
        echo "not_a_brick"
        echo "WARNING:get_brick_type:As it begins by a number this path should be considered 
    as a brick but is not : $brick_path" >&2
        return 1
    fi
    return 0
}

function convert_to_elementary_bricks_path { #< bricks_paths_list
    #> elementary_bricks_paths_list
    bricks_list="$1"
    return_code=0
    for brick in $bricks_list; do
        brick_type="$(get_brick_type "$brick")"
        if [ "$?" != 0 ]; then
            echo "ERROR:convert_to_elementary_bricks_path:get_brick_type:$brick" >&2
            return_code=1
        elif [ "$brick_type" == "super_brick" ]; then
            get_child_bricks "$brick"
            if [ $? != 0 ]; then
                return_code=1
            fi
        else
            echo "$brick"
        fi
    done
    return $return_code
}

function get_child_bricks { #< super_brick_path
    #> childs_bricks_paths_ordered_list
    brick_path="$1"
    get_all_bricks_paths | grep "^$brick_path/"
    return $?
}

function get_all_bricks_paths { #< nothing but read global ROOMS_LIST
    #> all_rooms_elementary_bricks_paths_ordered_list
    bricks_path_list="$(for room_path in $ROOMS_LIST ; do
        cd "$room_path"
        find . | grep "/[0-9]\+-[^/]*$" | grep -v '/[^0-9]' |
        sed "s|^\./|$room_path/|g" ; done)"
    for brick_path in $bricks_path_list ; do
        if [ "$(get_brick_type "$brick_path")" != "super_brick" ]; then
            echo "$brick_path"
        fi
    done
}

function get_all_bricks_names { #< nothing
    #> all_rooms_elementary_bricks_names_ordered_list
    for brick_path in $(get_all_bricks_paths); do
        get_brick_name "$brick_path"
    done
}

function display_bricks_in_right_order { #< brick_paths_list_to_display
    #> bricks_paths_ordered_list
    bricks_to_display="$1"
    get_all_bricks_path | while read brick ; do
        grep "^$brick$" <<<"$brick_to_display"
    done
    return 0
}

function get_brick_name { #< brick_path
    #> brick_name
    brick_path="$1"
    for room_path in $ROOMS_LIST ; do
        if grep -q "$room_path" <<<"$brick_path" ; then
            room_name="$(sed 's|^.*/\([^/]*\)$|\1|g' <<<"$room_path")"
            sed "s|$room_path|$room_name|g" <<<"$brick_path"
        fi
    done
}

function get_brick_path { #< brick_name
    #> brick_path
    brick_name="$1"
    room_name="$(sed 's|^\([^/]*\)/.*$|\1|g' <<<"$brick_name")"
    room_path="$(grep "/$room_name$" <<<"$ROOMS_LIST")"
    if [ -z "$room_path" ]; then # if $1 wasn't a brick name
        return 1
    fi
    parent_room_path="$(sed "s|/$room_name$||g" <<<"$room_path")"
    echo "$parent_room_path/$brick_name"
}

function get_bricks_paths_list { #< bricks_list # paths or names
    #> bricks_paths_list # doesn't check bricks validity
    output="$1"
    for room_path in $ROOMS_LIST; do
        room_name="$(sed 's|^.*/\([^/]*\)$|\1|g' <<<"$room_path")"
        output="$(sed "s|^$room_name|$room_path|g" <<<"$output")"
    done
    echo "$output"
}

function get_bricks_names_list { #< bricks_list # paths or names
    #> bricks_name_list # doesn't check bricks validity
    output="$1"
    for room_path in $ROOMS_LIST; do
        room_name="$(sed 's|^.*/\([^/]*\)$|\1|g' <<<"$room_path")"
        output="$(sed "s|^$room_path|$room_name|g" <<<"$output")"
    done
    echo "$output"
}

