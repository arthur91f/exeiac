#!/bin/bash
script_dir="$(cd "$(dirname "$0")"; pwd)"
exit_status=0

function get_tests_list {
    file="$1"
    python3 -c "import sys, yaml ; y=yaml.safe_load(sys.stdin.read()) ; print(y, end = \"\")" <"$file"
}

function get_tests_list_length {
    tests_list="$1"
    python3 -c "print(len($tests_list), end = \"\")"
}

function get_item {
    tests_list="$1"
    indice="$2"
    python3 -c "print($tests_list[$indice], end = \"\")"
}

function get_field {
    item="$1"
    field="$2"
    python3 -c "print($item[\"$field\"], end = \"\")"
}

function update_status_code {
    old_status="$1"
    new_status="$2"

    if [ "$old_status" -gt "$new_status" ]; then
        echo "$old_status"
    else
        echo "$new_status"
    fi
}

for file in $(find "$script_dir" -name '*.yml' | sort); do
    file_name="$(sed 's|^.*/\([^/]*\)\.yml$|\1|g' <<<"$file")"
    
    indice=0
    if ! tests_list="$(get_tests_list "$file")" ; then
        echo -e "\e[01;31mERROR\e[0;31m:test file not valid: $file : should be a yaml format\n\e[0;0m" >&2
        exit_status="$(update_status_code $exit_status 4)"
        continue
    fi
    
    if ! length="$(get_tests_list_length "$tests_list")" ; then
        echo -e "\e[01;31mERROR\e[0;31m:test file not valid: $file : should be a list\n\e[0;0m" >&2
        exit_status="$(update_status_code $exit_status 4)"
        continue
    fi

    for i in $(seq 0 $(($length-1))) ; do
        if ! item="$(get_item "$tests_list" $i )"; then
            echo -e "\e[01;31mERROR\e[0;31m:item $i not present: $tests_list\n\e[0;0m" >&2
            exit_status="$(update_status_code $exit_status 3)"
            continue
        fi
        title="$(get_field "$item" title 2>/dev/null)"
        cmd="$(get_field "$item" cmd 2>/dev/null)"
        if [ -z "$cmd" ]; then
            echo -e "\e[01;31mERROR\e[0;31m: test file $file_name not valid: item with title \"$title\" without cmd\n\e[0;0m" >&2
            exit_status="$(update_status_code $exit_status 3)"
            continue
        fi

        result_stdout="$(bash -c "$cmd")"
        result_status="$?"
        pass="true"

        if stdout="$(get_field "$item" stdout 2>/dev/null)" ; then
            if ! diff -qs <(echo "$stdout") <(echo "$result_stdout") >/dev/null ; then
                echo -e "\e[01;33mFAILED\e[0;33m:stdout_diff:$file_name: $title\e[0;0m"
                diff --color <(echo "$stdout") <(echo "$result_stdout")
                echo ""
                pass="false"
            fi
        fi

        if status="$(get_field "$item" status 2>/dev/null)" ; then
            if [ "$status" != "$result_status" ]; then
                echo -e "\e[01;33mFAILED\e[0;33m:status_diff:$file_name: $title"
                echo -e "  is $result_status instead of $status\e[0;0m"
                pass="false"
            fi
        fi

        if [ "$pass" == "true" ]; then
            echo -e "\e[01;32mPASSED\e[0;32m:*:$file_name: $title\e[0;0m"
        else
            exit_status="$(update_status_code $exit_status 2)"
        fi
        
    done

done

case $exit_status in
    0)
        echo "Every tests are valid ! Be careful anyway they aren't exhaustive" >&2
    ;;
    2)
        echo "Some tests have failed" >&2
    ;;
    3)
        echo "Some item in test file haven't a valid format" >&2
    ;;
    4)
        echo "A test file is not valid" >&2
    ;;
    *)
        echo "A non expacted error occurs: $exit_status" >&2
        exit_status=5
    ;;
esac

exit $exit_status