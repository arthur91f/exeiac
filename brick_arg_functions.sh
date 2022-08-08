#!/bin/bash
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
        echo "Warning: As it begins by a number this path should be considered 
    as a brick but is not : $brick_path" >&2
        return 1
    fi
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

function get_bricks_paths_list { #< bricks_names_list
    #> bricks_paths_list # doesn't check bricks validity
    output="$bricks_names_list"
    for room_path in $ROOMS_LIST; do
        room_name="$(sed 's|^.*/\([^/]*\)$|\1|g' <<<"$room_path")"
        output="$(sed "s|^$room_path|$room_name|g" <<<"$output")"
    done
    echo "$output"
}

function get_bricks_names_list { #< bricks_paths_list
    #> bricks_name_list # doesn't check bricks validity
    output="$bricks_paths_list"
    for room_path in $ROOMS_LIST; do
        room_name="$(sed 's|^.*/\([^/]*\)$|\1|g' <<<"$room_path")"
        output="$(sed "s|^$room_name|$room_path|g" <<<"$output")"
    done
    echo "$output"
}

