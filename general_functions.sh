#!/bin/bash
# Here are functions that aren't sticky to exeIaC

function get_arg {
    # get_arg --TYPE=ARG_SEARCHED ARG1 ARG2 AR3 ...
    type="$(sed 's|^-*\([a-zA-Z0-9_-]*\)=.*$|\1|g' <<<"$1")"
    arg_name="$(sed 's|^-*[a-zA-Z0-9_-]*=\(.*\)$|\1|g' <<<"$1" |
        sed 's|^-\{1,2\}||g')"
    shift
    previous_arg=""
    case "$type" in
    boolean|b)
        for arg in "$@" ; do
            if [ "$arg_name" == "$(sed 's|^-\{1,2\}||g' <<<"$arg")" ]; then
                return 0
            elif grep -q "^-\{1,2\}$arg_name=true" <<<"$arg"; then
                return 0
            elif grep -q "^-\{1,2\}$arg_name=false" <<<"$arg"; then
                return 1
            fi
        done
        return 1
    ;;
    string|s)
        for arg in "$@" ; do
            if grep -q "^-\{1,2\}$arg_name=.*" <<<"$arg"; then
                sed "s|^-\{1,2\}$arg_name=\(.*\)$|\1|g" <<<"$arg"
                return 0
            elif grep -q "^-\{1,2\}$arg_name$" <<<"$arg"; then
                previous_arg="$arg_name"
            elif [ "$previous_arg" == "$arg_name" ]; then
                echo "$arg"
                return 0
            fi
        done
        return 1
    ;;
    *)
        echo "ERROR:get_arg:bad_argument:$type" >&2
        return 2
    ;;
    esac 
}

function get_arg_in_string {
    # string="quiet,string=\"I'm good, and you ?\",output='filename with space',type=file"
    # get_arg_in_string "$string" quiet ; echo $?
        # true
        # 0
    # get_arg_in_string "$string" string ; echo $?
        # I'm good, and you ?
        # 0
    # get_arg_in_string "$string" test-arg-not-present ; echo $?
        #
        # 1
    arg_searched="$1"
    args_list_string=",$2,"

    if grep -q ",$arg_searched.*," <<<"$args_list_string"; then
        if grep -q ",$arg_searched," <<<"$args_list_string"; then
            echo "true"
        elif grep -q ",$arg_searched='.*'," <<<"$args_list_string"; then
            sed "s/^.*,$arg_searched='\(.*\)',.*$/\1/g" <<<"$args_list_string"
        elif grep -q ",$arg_searched=\".*\"," <<<"$args_list_string"; then
            sed "s/^.*,$arg_searched=\"\(.*\)\",.*$/\1/g" <<<"$args_list_string"
        elif grep -q ",$arg_searched=.*," <<<"$args_list_string"; then
            sed "s/^.*,$arg_searched=\([^,]*\),.*$/\1/g" <<<"$args_list_string"
        fi
        return 0
    else
        return 1
    fi
}

function display_line_after_match {
    text="$1"
    regex="$2"
    match_line="$(echo "$text" | grep -n "$regex" | cut -d: -f1)"
    if [ -z "$match_line" ]; then
        return 1
    else
        echo "$text" | sed -n "$match_line,\$p"
    fi
}

function display_line_before_match {
    text="$1"
    regex="$2"
    match_line="$(echo "$text" | grep -n "$regex" | cut -d: -f1)"
    if [ -z "$match_line" ]; then
        return 1
    else
        echo "$text" | sed -n "1,${match_line}p"
    fi
}

function get_absolute_path {
    path="$1"
    if [ -e "$path" ]; then
        echo "$(cd "$(dirname "$path")"; pwd)/$(basename "$path")"
    else
        echo "Error path not exist: $path"
        return 1
    fi
}

function merge_string_on_new_line {
	if [ -n "$1" ]; then
		echo "$1"
	fi
	if [ -n "$2" ]; then
		echo "$2"
	fi
}

