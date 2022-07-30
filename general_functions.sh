#!/bin/bash
# Here are functions that aren't sticky to exeIaC

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

