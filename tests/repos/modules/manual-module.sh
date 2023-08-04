#!/bin/bash
ACTION="$1"
ALL_ARGS="$@"
CURRENT_PATH="$(pwd)"
#FILES_LIST="$EXEIAC_files_list" until we support value data
FILES_LIST="users.yml
groups.yml"

function describe_module_for_exeiac {
    echo '{
    "init": {
        "behaviour": "standard"
    },
    "plan": {
        "behaviour": "plan"
    },
    "lay": {
        "behaviour": "lay"
    },
    "remove": {
        "behaviour": "remove"
    },
    "output": {
        "behaviour": "output"
    },
    "validate_code": {
        "behaviour": "standard"
    }
}'
    exit 0
}

function show_implemented_actions {
    grep "^function " $0 |
        sed 's|^function \([^ ]*\) .*$|\1|g' |
        grep -v "internal"
    exit 0
}

function init {
    echo "execute init: install tools and deps for $PWD"

    status_code=0
    python3 --version 2>/dev/null
    if [ "$?" != 0 ]; then
        echo "manual-module:init: python3 not installed"
        return 22
    fi
    if which jq >/dev/null ; then
        echo "manual-module:init: jq installed"
    else
        echo "manual-module:init: jq not installed"
        echo "  https://stedolan.github.io/jq/download/"
        return 23
    fi
}

function internal_ask_confirmation {
    echo -e "\033[1m$1\033[0m\033[3m(only yes accepted): \033[0m"
    read answer
    if [ "$answer" == "yes" ]; then
        return 0
    else
        return 1
    fi
}

function internal_get_file_basename {
    sed -e 's/.json$//g' -e 's/\.yml$//g' -e 's/\.yaml$//g' <<<"$1"
}

function plan {
    if ! validate_code; then
        return 21
    fi
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        echo "manual-module:plan: as the lay is a manual and the non-interactive mode is enabled, can't determine if there is a drift"
        return 3
    else
        output
        if internal_ask_confirmation "Is it correct ?" ; then
            return 0
        else
            return 2
        fi
    fi
}

function lay {
    if ! validate_code; then
        return 21
    fi
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        echo "manual-module:lay: as the lay is a manual and the non-interactive mode is enabled, assume everything is ok"
    else
        output
        if internal_ask_confirmation "Do you want to change anything ?" ; then
            for file in $FILES_LIST ; do
                file_basename="$(internal_get_file_basename "$file")"
                if internal_ask_confirmation "Do you want to modify $file_basename ?" ; then
                    nano $file
                fi
            done
        fi
    fi 
}

function remove {
    for file in $FILES_LIST ; do
        file_basename="$(internal_get_file_basename <<<"$file")"
        echo "" > $file
        echo "$file_basename"
    done
}

function output {
    (echo "{" ; for file in $FILES_LIST ; do
        if grep -q "\.json$" <<<"$file"; then
            file_basename="$(internal_get_file_basename "$file")"
            echo -n "\"$file_basename\": {$(cat $file)},"
        elif grep -Eq "\.(yml|yaml)" <<<"$file"; then
            file_basename="$(internal_get_file_basename "$file")"
            json="$(python3 -c 'import sys, yaml, json; y=yaml.safe_load(sys.stdin.read()); print(json.dumps(y))' < $file)"
            echo -n "\"$file_basename\": $json,"
        else
            echo "manual-module:output: format not recognized or supported for $file" >&2
            return 1
        fi
    done | sed 's/,$//g'; echo "}") | jq .
}

function validate_code {
    list="$(for file in $FILES_LIST ; do
        internal_get_file_basename "$file"
    done | sort | uniq -c | sed '/^ *1/d')"
    
    if [ -n "$list" ]; then
        echo "Some files are in double : "
        list="$(echo "$list" | sed 's/^ *[0-9]* //g')"
        for item in $list ; do
            grep "$item" <<<"$FILES_LIST"
        done
        return 21
    fi
}

if grep -q "^$ACTION$" <(show_implemented_actions) ; then
    $ACTION
    exitcode="$?"
else
    echo "action not implemented: $ACTION"
    exitcode=3
fi

exit "$exitcode"
