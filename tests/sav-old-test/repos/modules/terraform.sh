#!/bin/bash
ACTION="$1"
ALL_ARGS="$@"
CURRENT_PATH="$(pwd)"
PLAN_PATH="$CURRENT_PATH/.exeiac_module.tfplan"

function describe_module_for_exeiac {
    echo '{
    "init": {
        "behaviour": "standard"
    },
    "plan": {
        "behaviour": "plan",
        "status_code_fail": "1,3-255",
        "events": {
            "exeiac_plan_no_drift": {
                "type": "status_code",
                "status_code": "0"
            },
            "exeiac_plan_drift": {
                "type": "status_code",
                "status_code": "2"
            }
        }
    },
    "lay": {
        "behaviour": "lay",
        "status_code_fail": "1,3-255",
        "events": {
            "nothing_todo": {
                "type": "status_code",
                "status_code": "0"
            },
            "drift": {
                "type": "status_code",
                "status_code": "2"
            },
            "recreated_resources": {
                "type": "file",
                "path": "./.exeiac_events"
            }
        }
    },
    "remove": {
        "behaviour": "remove"
    },
    "output": {
        "behaviour": "output"
    },
    "clean": {
        "behaviour": "clean"
    }
}'
    exit 0
}

function internal_quit {
    text="$1"
    exit_code="$2"

    if [ -z "$text" ]; then
        true
    elif [ "$(wc -l <<<"$text")" == 1 ]; then
        echo "$text" >&2
    else
        indent="terraform:"
        while read line ; do
            echo "$indent$line" >&2
            indent="    "
        done <<<"$text"
    fi

    exit $exit_code

}

function init {
    echo "execute init: install tools and deps for $PWD"

    status_code=0
    terraform version 2>/dev/null
    if [ "$?" != 0 ]; then
        echo "terraform:init: terraform not installed"
        echo "  https://developer.hashicorp.com/terraform/downloads"
        status_code=1
    fi
    if which jq >/dev/null ; then
        echo "terraform:init: jq installed"
    else
        echo "terraform:init: jq not installed"
        echo "  https://stedolan.github.io/jq/download/"
        status_code=$(($status_code+2))
    fi

    terraform init
    err=$?
    if [ "$err" != 0 ]; then
        return $err
    elif [ "$status_code" != 0 ]; then
        status_code=$(($status_code+20))
        return $status_code
    else
        return 0
    fi
}

function plan {
    terraform plan -detailed-exitcode
    return $?
}

function internal_create_event {
    cat "$PLAN_PATH.txt" |
        grep -E '^(\-/\+|\+|\-) resource "[^"]*" "[^"]*" {' | 
        sed 's|^\(.*\) resource "\([^"]*\)" "\([^"]*\)" {|\1 \2.\3|g' > "$CURRENT_PATH/.exeiac_events"
}

function lay {
    terraform plan -out="$PLAN_PATH" > "$PLAN_PATH.txt"
    plan_status_code="$?"

    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        terraform apply -auto-approve "$PLAN_PATH"
        status_code="$?"
    else
        terraform apply "$PLAN_PATH"
        status_code="$?"
    fi
    if [ "$status_code" != 0 ]; then
        echo "terraform apply failed with code $status_code" >&2
        return 11
    fi

    internal_create_event
    if [ "$?" != 0 ]; then
        echo "terraform module failed to create events: $status_code" >&2
        return 12
    fi

    return "$plan_status_code"
}

function remove {
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        terraform destroy -auto-approve
        return $?
    else
        terraform destroy
        return $?
    fi
}

function output {
    json="$(terraform output -json)"
    err="$?"
    if [ "$err" != 0 ]; then
        echo "{}"
        return $err
    fi
    jq 'map_values(.value)' <<<"$json"
    err="$?"
    if [ "$err" != 0 ]; then
        echo "{}"
        return $err
    fi
    return 0
}

function clean {
    if [ -e ./.terraform ]; then
        rm -r ./.terraform || 
        internal_quit "clean:error when 'rm -r .terraform'" 21
    fi
    if [ -e ./.terraform.lock.hcl ]; then
        rm ./.terraform.lock.hcl ||
        internal_quit "clean:error when 'rm -r .terraform.lock.hcl'" 22
    fi
    if [ "$(output)" == "{}" ] && [ -e /terraform.tfstate ] ; then
        rm ./terraform.tfstate || 
            internal_quit "clean:error when 'rm terraform.tfstate'" 23
    fi
    if [ -f "./.exeiac_events" ]; then
        rm "./.exeiac_events"
    fi
    return 0
}

if [ "$ACTION" == "describe_module_for_exeiac" ]; then
    describe_module_for_exeiac
elif grep -q "^$ACTION$" <(describe_module_for_exeiac | jq -r 'keys | .[]') ; then
    $ACTION
    exitcode="$?"
else
    echo "action not implemented: $ACTION"
    exitcode=21
fi

exit "$exitcode"
