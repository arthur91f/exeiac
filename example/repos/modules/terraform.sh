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

function lay {
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        terraform apply -auto-approve
        return $?
    else
        terraform apply
        return $?
    fi
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
    echo "$json" | jq 'map_values(.value)'
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
    return 0
}

if grep -q "^$ACTION$" <(show_implemented_actions) ; then
    $ACTION
    exitcode="$?"
else
    echo "action not implemented: $ACTION"
    exitcode=21
fi

exit "$exitcode"
