#!/bin/bash
copy_function show_dependencies default_show_dependencies

function is_brick_using_this_module {
    brick_path="$(pwd)"
    if [ -d "$brick_path" ] && ls -1 "$brick_path" | grep -q "\.tf$" &&
        ! grep -q "^[0-9]\+-[^/]*$" <<<"$(ls -1 "$brick_path")"; then
        return 0
    else
        return 1
    fi
}

function install {
    terraform -v
    # os_type="linux_amd64"
    # install_path=/usr/local/bin
    # version="$(curl -s https://releases.hashicorp.com/terraform/ |
    #  grep '<a href="/terraform/.*/">terraform_[0-9.]*</a>' |
    #  sed 's|^ *<a href="/terraform/[0-9.]*/">terraform_\([0-9.]*\)</a>$|\1|g' |
    #  head -n1)"
    # download_url="https://releases.hashicorp.com/terraform/$version/ terraform_${version}_${os_type}.zip"
    # sudo wget "$download_url" -O "/tmp/terraform.zip"
    # sudo unzip /tmp/terraform.zip -d /tmp
    # sudo rm /tmp/terraform.zip
    # sudo mv /tmp/terraform "$install_path/terraform_$version"
    # sudo chmod 755 "$install_path/terraform_$version"
    # [ -f "$install_path/terraform" ] && rm "$install_path/terraform"
    # ln -s "$install_path/terraform_$version" "$install_path/terraform"
}

function init {
    terraform init $MODULE_OPTS
}

function plan {
    terraform plan -detailed-exitcode $MODULE_OPTS
}

function apply {
    if grep -q "non-interactive" <<<"$EXEIAC_OPTS"; then
        terraform apply -auto-approve $MODULE_OPTS
    else
        terraform apply $MODULE_OPTS
    fi
}

function output {
    if grep -q "dry-terraform-output" <<<"$EXEIAC_OPTS"; then
        terraform output $MODULE_OPTS
    else
        jq_filter=".\"$MODULE_OPTS\""
        output="$(terraform output -json | jq 'map_values(.value)' | jq
            "$jq_filter")"
        if grep -q '^".*"$'<<<"$output" ; then
            echo "$output" | sed -e 's/^"//g' -e 's/"$//g'
        else
            echo "$output"
        fi
    fi
}

function destroy {
    if grep -q "non-interactive" <<<"$EXEIAC_OPTS"; then
        terraform destroy -auto-approve $MODULE_OPTS
    else
        terraform destroy $MODULE_OPTS
    fi
}

function validate {
    terraform validate $MODULE_OPTS
}

function fmt {
    terraform fmt $MODULE_OPTS
}

function show_dependencies {
    comment_dependencies="$(default_show_dependencies)"
    gcs_backend_dependencies="$(cat "$brick_path"/*.tf |
        sed -n '/^data "terraform_remote_state" ".*" {/,/^}/ p' |
        sed -n '/  config = {/,/  }/ p' |
        grep "prefix *=" | 
        sed 's|^ *prefix *= *"\(.*\)".*$|\1|g')"

    echo "$(echo "$comment_dependencies" ;
        echo "$gcs_backend_dependencies")" |
        sed '/^$/d' | sort | uniq
}

function help {
    echo "The help of this brick haven't been overloaded so it runs the normal way."
    echo "execiac BRICK_PATH ACTION [OPTIONS]"
    echo "init: run terraform init"
    echo "plan: run terraform plan -detailed-exitcode"
    echo "apply: run terraform apply"
    echo "output [OUTPUT_FIELD] [--exeiac-opts=dry-terraform-output]: terraform output
    in json"
    echo "destroy: run terraform destroy"
    echo "validate: run terraform validate"
    echo "fmt: rewrite file to pass linter"
    echo "show_dependencies: search comment #TAG:depends:BRICK_NAME"
    echo "install: install terraform in /usr/local/bin"
    echo ""
    echo "special options:"
    echo "--exeiac-opts=dry-terraform-output: permit to only run terraform output
    althought use json output"
    echo "--exeiac-opts=non-interactive: add -auto-approve opts to terraform apply, note that it isn't necessary to put --exeiac-opts= before it's just for other modules compatibilities"
}

