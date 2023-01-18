#!/bin/bash
ACTION="$1"
ALL_ARGS="$@"
CURRENT_PATH="$(pwd)"
OUTPUT_FILE="$CURRENT_PATH/output.json"

function show_implemented_actions {
    grep "^function " $0 |
        sed 's|^function \([^ ]*\) .*$|\1|g' |
        grep -v "internal"
    exit 0
}

function ask_confirmation {
    echo -en "\033[1m$1\033[0m\033[3m (only yes accepted): \033[0m"
    read answer
    if [ "$answer" == "yes" ]; then
        return 0
    else
        return 1
    fi
}

function init {
    if which jq >/dev/null ; then
        echo "jq installed: ok"
    else
        echo "jq not installed"
        return 3
    fi
}

function plan {
    echo "plan: as the lay is a manual step we assume everything is ok"
    return 0
}

function lay {
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS"; then
        echo "lay: non-interactive -> assume everything is ok"
        return 0
    fi
    
    echo "Do it by hand"
    echo "Create account on cloud provider"
    for cloud_provider in $(jq -r '.cloud_providers | keys | .[]' "$OUTPUT_FILE") ; do
        url="$(jq -r ".cloud_providers.\"$cloud_provider\".\"signup_url\"" "$OUTPUT_FILE")"
        org=$(jq -r ".cloud_providers.\"$cloud_provider\".\"organisation\"" "$OUTPUT_FILE")
        echo "- Create an account \"$org\" on $cloud_provider ($url)"
    done
    echo -e "\nCreate project on cloud provider"
    for project in $(jq -r '.projects | keys | .[]' "$OUTPUT_FILE") ; do
        cloud_provider="$(jq -r ".projects.\"$project\".cloud_provider" "$OUTPUT_FILE")"
        echo "- Create project $project on $cloud_provider"
    done
    if ask_confirmation "enter yes when it's done" ; then
        return 0
    else
        return 1
    fi
}

function remove {
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS"; then
        echo "remove: non-interactive -> assume everything is ok"
        return 0
    fi
    
    echo "Do it by hand"
    echo "- Close your account and remove your creditcard on localtest"
    if ask_confirmation "enter yes when it's done" ; then
        return 0
    else
        return 1
    fi
}

function output {
    jq '.' "$OUTPUT_FILE"
}

if grep -q "^$ACTION$" <(show_implemented_actions) ; then
    $ACTION
    exitcode="$?"
else
    echo "action not implemented: $ACTION"
    exitcode=3
fi

exit "$exitcode"
