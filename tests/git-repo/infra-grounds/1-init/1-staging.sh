#!/bin/bash
state_file="$(pwd)/staging_state.cyphered"

function init {
    if which jq >/dev/null; then
        echo "jq not installed" >&2
    fi
}

function plan {
    if [ -e "$state_file" ]; then
        echo "Nothing to do"
        return 0
    else
        return 2
        echo "Display manual action to do for prod"
        echo "Create file: $state_file"
    fi
    return 1
}

function apply {
    plan
    if plan ; then
        echo "Apply ok, nothing to do"
        return 0
    else
        if ! get_arg --boolean=non-interactive "${OPTS[@]}"; then
            echo -n "enter \"yes\" when it's done: "
            ask_confirmation
            if plan >/dev/null; then
                return 0
            else
                echo "apply as failed"
                return 1
            fi
        fi
        return 1
    fi
}

function output {
    if [ -e "$state_file" ]; then
        cat "$state_file" | jq
    fi
}

function destroy {
    if plan >/dev/null; then
        echo "Delete cloud project $(cat "$state_file" | jq ".project_name")"
        if ! get_arg --boolean=non-interactive "${OPTS[@]}"; then
            echo -n "enter \"yes\" when it's done: "
            ask_confirmation
            if plan >/dev/null; then
                echo "Destroy as failed"
                return 1
            else
                return 0
            fi
        else
            return 1
        fi
    else
        echo "Destroy: nothing to do"
        return 0
    fi
}

