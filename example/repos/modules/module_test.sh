#!/bin/bash
ACTION="$1"
ALL_ARGS="$@"
CURRENT_PATH="$(pwd)"

function show_implemented_actions {
    grep "^function " $0 |
        sed 's|^function \([^ ]*\) .*$|\1|g' |
        grep -v "internal"
    exit 0
}

function init {
    echo "execute init: install tools and deps for $PWD"
}

function internal_interactive {
    echo "execute $ACTION: for test you can choose what's happen"
    read -p "choose: ok|drift|fail: " answer
    if [ "$answer" == "ok" ]; then
        echo "no drift"
        return 0
    elif [ "$answer" == "drift" ]; then
        echo "some diff have been found"
        return 1
    else
        echo "fail"
        return $(( $RANDOM % 255 + 2 ))
    fi
}

function plan {
    # for the test the result will depends of the presence of word in path
    if grep -q "drift$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: some diff have been found, need to be layed"
        return 1
    elif grep -q "fail$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: the $ACTION have failed"
        echo "remove the ending \"fail\" of the brick name" >&2
        return $(( $RANDOM % 255 + 2 ))
    elif grep -q ".*--non-interactive" <<<"$ALL_ARGS"; then
        echo "execute $ACTION: ok no drift"
        return 0
    elif grep -q "interactive" <<<"$BRICK_PATH"; then
        internal_interactive
        return $?
    else
        echo echo "execute plan: ok no drift"
        return 0
    fi
}

function lay {
    # for the test the result will depends of the presence of word in path
    if grep -q "drift$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: some diff have been found, need to be layed"
        echo "  - null_resource: test"
        echo "Do you want to lay (only \"yes\" accepted):"
        read answer
        if [ "$answer" == "yes" ]; then
            echo "lay ok"
            return 1
        else
            return 2
        fi
    elif grep -q "fail$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: the $ACTION have failed"
        echo "remove the ending \"fail\" of the brick name" >&2
        return $(( $RANDOM % 255 + 2 ))
    elif grep -q ".* --non-interactive" <<<"$ALL_ARGS"; then
        echo "execute $ACTION: ok no drift"
        return 0
    elif grep -q "interactive" <<<"$BRICK_PATH"; then
        internal_interactive
        return $?
    else
        echo echo "execute plan: ok no drift"
        return 0
    fi
}

function remove {
    lay
}

function output {
    # for test we display all outputs here it's simpler
    echo '{
    "production": {
        "name": "production",
        "project": "prod-120822",
        "credntials": {
            "user": "ci-prod",
            "key": "-----BEGIN PRIVATE KEY-----\nproductionfjh6qnrS4z3qPOrOthu9hVJiM9sXPISuIMWCQUVDhA7cKwuZ3ErN3Mue2GBt\nKYuINui3SbsyfuAJlF/3FcUVfbwbioRWMKEs4rTPeU9Nz06Ipj2wZnRUpQ0njgptzxyrzH\nNblRVU1Skj9xNTzNbj02bWoOYSOgq3YIF+901gbQuZdE8JPxcKZMUXSxTvhgq5zA5bpKrj\nZm3Gja0rYjNM6NEFCRfVW4Fg==\n-----END PRIVATE KEY-----"
        }
    },
    "staging": {
        "name": "staging",
        "project": "stage-120822",
        "credntials": {
            "user": "ci-staging",
            "key": "-----BEGIN PRIVATE KEY-----\nstagingfjh6qnrS4z3qPOrOthu9hVJiM9sXPISuIMWCQUVDhA7cKwuZ3ErN3Mue2GBt\nKYuINui3SbsyfuAJlF/3FcUVfbwbioRWMKEs4rTPeU9Nz06Ipj2wZnRUpQ0njgptzxyrzH\nNblRVU1Skj9xNTzNbj02bWoOYSOgq3YIF+901gbQuZdE8JPxcKZMUXSxTvhgq5zA5bpKrj\nZm3Gja0rYjNM6NEFCRfVW4Fg==\n-----END PRIVATE KEY-----"
        }
    },
    "monitoring": {
        "name": "monitoring",
        "project": "monit-120822",
        "credntials": {
            "user": "ci-monitoring",
            "key": "-----BEGIN PRIVATE KEY-----\nmonitoringfjh6qnrS4z3qPOrOthu9hVJiM9sXPISuIMWCQUVDhA7cKwuZ3ErN3Mue2GBt\nKYuINui3SbsyfuAJlF/3FcUVfbwbioRWMKEs4rTPeU9Nz06Ipj2wZnRUpQ0njgptzxyrzH\nNblRVU1Skj9xNTzNbj02bWoOYSOgq3YIF+901gbQuZdE8JPxcKZMUXSxTvhgq5zA5bpKrj\nZm3Gja0rYjNM6NEFCRfVW4Fg==\n-----END PRIVATE KEY-----"
        }
    },
    "domain_name": {
        "internal": "internal.mycompany.co",
        "private": "priv.mycompany.co",
        "public": "mycompany.com"
    },
    "network": {
        "ip_range": "10.0.0.0/20",
        "network_id": "myaccount/myproject/network/123456-7890"
    }
}'
}

$ACTION
EXIT_CODE="$?"

exit "$EXIT_CODE"
