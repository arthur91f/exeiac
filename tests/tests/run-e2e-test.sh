#!/bin/bash
script_dir="$(cd "$(dirname "$0")"; pwd)"

function get_field {
    field="$1"
    item="$2"
    value=""
    if grep -Eq "^(-| ) $field: *\| *$" <<<"$item"; then
        match_line="$(($(echo "$item" | grep -En "^(-| ) $field: *\| *$" | cut -d: -f1 | head -n1)+1))"
        text_after_match="$(sed -n "$match_line,\$p" <<<"$item")"
        if grep -q "^  [^ ]" <<<"$text_after_match"; then
            match_line="$(($(grep -n "^  [^ ]" <<<"$text_after_match" | cut -d: -f1 | head -n1)-1))"
            text="$(sed -n "1,${match_line}p" <<<"$text_after_match")"
            value="$(sed 's/^    //g' <<<"$text")"
        else
            text="$text_after_match"
        fi
        value="$(sed 's/^    //g' <<<"$text")"
    elif grep -q "^- $field:" <<<"$item"; then
        value="$(grep "^- $field:" <<<"$item" | 
            sed -e "s/^- $field://g" -e 's/ *$//g' -e 's/^ *//g' -e 's/ *#.*$//g')"
    elif grep -q "^  $field:" <<<"$item"; then
        value="$(grep "^  $field:" <<<"$item" | 
            sed -e "s/^  $field://g" -e 's/ *$//g' -e 's/^ *//g' -e 's/ *#.*$//g')"
    else
        return 1
    fi
    echo "$value"
}

function get_item {
    content="$1"
    start_line="$2"
    if [ "$start_line" == 0 ]; then
        text_after_match="$content"
    else
        text_after_match="$(sed -n "$start_line,\$p" <<<"$content")"
    fi

    end_line="$(($(grep -n "^- [^ ]" <<<"$text_after_match"| cut -d: -f1 | head -n2 | tail -n1)-1))"
    if [ "$end_line" == 1 ]; then
        echo "$text_after_match" | tail -n+2
    else
        if [ "$start_line" == 0 ]; then
            sed -n "1,${end_line}p" <<<"$text_after_match"
        else
            sed -n "1,${end_line}p" <<<"$text_after_match" | tail -n+2
        fi
    fi
}

exit_status=0

for file in $(find "$script_dir" -name '*.yml' | sort); do
    file_name="$(sed 's|^.*/\([^/]*\)\.yml$|\1|g' <<<"$file")"
    
    # remove commented line
    content="$(sed '/^ *#/d' "$file")"
    for start_line in $(grep -n "^- " <<<"$content" | cut -d: -f1); do
        item="$(get_item "$content" "$(($start_line - 1))")"
        
        title="$(get_field title "$item")"
        cmd="$(get_field cmd "$item")"
        if [ -z "$cmd" ]; then
            echo "ERROR: test file $file_name not valid: item with title \"$title\" without cmd" >&2
            exit 3
        fi

        result_stdout="$(bash -c "$cmd")"
        result_status="$?"
        pass="true"

        if stdout="$(get_field stdout "$item")" ; then
            if ! diff -qs <(echo "$stdout") <(echo "$result_stdout") >/dev/null ; then
                echo "FAILED:stdout_diff:$file_name: $title"
                diff --color <(echo "$stdout") <(echo "$result_stdout")
                pass="false"
            fi
        fi

        if status="$(get_field status "$item")" ; then
            if [ "$status" != "$result_status" ]; then
                echo "FAILED:status_diff:$file_name: $title"
                echo "  is $result_status instead of $status"
                pass="false"
            fi
        fi

        if [ "$pass" == "true" ]; then
            echo "PASSED:*:$file_name: $title"
        else
            exit_status=2
        fi
                
    done

done

exit $exit_status
