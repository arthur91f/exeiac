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
        return 1
    fi
    if which jq >/dev/null ; then
        echo "terraform:init: jq installed"
    else
        echo "terraform:init: jq not installed"
        echo "  https://stedolan.github.io/jq/download/"
        status_code=1
    fi

    terraform init
}

#function internal_ask_confirmation {
#     echo -e "\033[1m$1\033[0m\033[3m(only yes accepted): \033[0m"
#     read answer
#     if [ "$answer" == "yes" ]; then
#         return 0
#     else
#         return 1
#     fi
# }

function plan {
    terraform plan -detailed-exitcode
    case "$?" in
    1)
        return 2 ;;
    2)
        return 1 ;;
    *)
        return $? ;;
    esac
}

function lay {
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        terraform apply -auto-approve
    else
        terraform apply
    fi
}

function remove {
    if grep -q ".*--non-interactive" <<<"$ALL_ARGS" ; then
        terraform destroy -auto-approve
    else
        terraform destroy
    fi
}

function output {
    terraform output -json | jq 'map_values(.value)'
}

function clean {
    rm -r ./.terraform || 
        internal_quit "clean:error when 'rm -r .terraform'" 2
    rm ./.terraform.lock.hcl ||
        internal_quit "clean:error when 'rm -r .terraform.lock.hcl'" 2
    if [ "$(output)" == "{}" ]; then
        rm ./terraform.tfstate || 
            internal_quit "clean:error when 'rm terraform.tfstate'" 2
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
