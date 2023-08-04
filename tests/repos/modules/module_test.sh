#!/usr/bin/env bash
ACTION="$1"
ALL_ARGS="$@"
CURRENT_PATH="$(pwd)"

function internal_ask_confirmation {
    echo -e "\033[1m$1\033[0m\033[3m(only yes accepted): \033[0m"
    read answer
    if [ "$answer" == "yes" ]; then
        return 0
    else
        return 1
    fi
}

function internal_get_toCreate {
    (echo -n "{"
    compgen -v | grep "^EXEIAC_TEST_" | while read varname ; do
        field="$(sed 's/^EXEIAC_TEST_//g' <<<"$varname")"
        value="$(eval echo \"\$$varname\")"
        if grep -q "^{" <<<"$value" ; then
            echo -n "\"from_$field\": $value,"
        else
            echo -n "\"from_$field\": \"$value\","
        fi
    done | sed 's/,$//g' ; echo "}") | jq --sort-keys .
}

function internal_get_created {
    if [ -f CREATED_this ]; then
        cat CREATED_this | jq --sort-keys .
    else
        echo "{}"
    fi
}

function internal_display_diff {
    diff --color -y <(echo "$1" | jq --sort-keys .) <(echo "$2" | jq --sort-keys .)
    err=$?
    case "$err" in
        1) return 2 ;;
        2) return 1 ;;
        *) return $err ;;
    esac
}

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
    if which jq >/dev/null ; then
        echo "test-module:init: jq installed"
        return 0
    else
        echo "test-module:init: jq not installed"
        echo "  https://stedolan.github.io/jq/download/"
        return 21
    fi
}

function plan {
    internal_display_diff "$(internal_get_created)" "$(internal_get_toCreate)"
    return $?
}

function lay {
    to_create="$(internal_get_toCreate)"
    internal_display_diff "$(internal_get_created)" "$to_create"
    if [ "$?" == 0 ]; then
        return 0
    fi
    
    if ! grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        if ! internal_ask_confirmation "Do you want to continue ? " ; then
            return 21
        fi
    fi
    echo "$to_create" > CREATED_this

}

function remove {
    if [ "$?" == 0 ]; then
        return 0
    fi
    internal_display_diff "$(internal_get_created)" "{}"
    echo "{}" > CREATED_this
}

function output {
    internal_get_created
}

$ACTION
EXIT_CODE="$?"

exit "$EXIT_CODE"
