#!/bin/bash
# Here are functions that aren't sticky to exeIaC

function get_arg { #< =(-boolean=|-string=)arg_name ARG1 ARG2...
    #> ? return 0 if the arg is found
    #> ? if $1="-boolean=arg_name" can return false if arg_name=false
    #> if $1="-string=arg_name" return arg_value

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
                return 2
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

function display_line_after_match { #< multiline_text_of_uniq_line regex
    #> ? #> truncated_multiline_text_beginning_with_match
    text="$1"
    regex="$2"
    match_line="$(echo "$text" | grep -n "$regex" | cut -d: -f1)"
    if [ -z "$match_line" ]; then
        return 1
    else
        echo "$text" | sed -n "$match_line,\$p"
    fi
}

function display_line_before_match { #< multiline_text_of_uniq_line regex
    #> ? #> truncated_multiline_text_ending_with_match
    text="$1"
    regex="$2"
    match_line="$(echo "$text" | grep -n "$regex" | cut -d: -f1)"
    if [ -z "$match_line" ]; then
        return 1
    else
        echo "$text" | sed -n "1,${match_line}p"
    fi
}

function get_absolute_path { #< relative_path_or_absolute_path
    #> absolute_path
    path="$1"
    if [ -e "$path" ]; then
        echo "$(cd "$(dirname "$path")"; pwd)/$(basename "$path")"
    else
        echo "Error path not exist: $path"
        return 1
    fi
}

function merge_string_on_new_line { #< string1 string2 # can be empty or multiline
    #> empty_or_multiline_string
	if [ -n "$1" ]; then
		echo "$1"
	fi
	if [ -n "$2" ]; then
		echo "$2"
	fi
}

function soft_exit { #< return_code error_message
    #> ? #>2 error_message # simply exit
    return_code="$1"
    cd "$INITIAL_CURRENT_PATH"
    if [ -n "$2" ]; then
        echo "$2" >&2
    fi
    return "$return_code"
}

