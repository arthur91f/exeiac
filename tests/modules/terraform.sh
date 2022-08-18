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
    return_code=0
    if which terraform >/dev/null ; then
        echo "terraform installed: ok"
    else
        echo "ERROR:terraform is not installed"
        return_code=1
    fi

    if which jq >/dev/null ; then
        echo "jq installed: ok"
    else
        echo "ERROR:jq is not installed"
        return_code=1
    fi
    return $return_code
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
    # CHECK INSTALL
    if which terraform >/dev/null ; then
        echo "terraform installed: ok"
    else
        echo "ERROR:terraform is not installed"
        return 1
    fi

    if which jq >/dev/null ; then
        echo "jq installed: ok"
    else
        echo "ERROR:jq is not installed"
        return 1
    fi

    # INIT
    terraform init
    return $?
}

function get_env {
    envs_list='(production|staging|monitoring)'
    brick_name="$(get_brick_name "$(pwd)")"
    if brick_name="$(get_brick_name "$(pwd)")"; then
        if env_dir="$(egrep -o "/[0-9]-$envs_list/" <<<"$brick_name")"; then
            sed -e 's|^/[0-9]-||g' -e 's|/$||g' <<<"$env_dir"
            return 0
        else
            return 1 
        fi
    else
        return 1
    fi
}

function get_vars_file_content {
    ## variable environment ##
    env="$(get_env)"
    init_brick_path="$(get_brick_path "infra-grounds/1-init")"
    vars_file_cyphered="$init_brick_path/${env}_state.cyphered"
    cat "$vars_file_cyphered"
    
    ## variable state_tag ##
    state_tag="$(get_brick_sanitized_name "$(get_brick_name "$(pwd)")")"
    echo "state_tag = \"$state_tag\""

    ## variable rooms_paths_list ##
    echo "rooms_paths_list = {"
    echo $room_paths_list | sed 's|^\(.*/\([^/]*\)\)$|  \2 = "\1"|g'
    echo "}"
}

function plan {
    state_tag="$(get_brick_sanitized_name "$(get_brick_name "$(pwd)")")"
    terraform plan -detailed-exitcode -var-file=<(get_vars_file_content)
    return $?
}

function apply {
    state_tag="$(get_brick_sanitized_name "$(get_brick_name "$(pwd)")")"
    if get_arg --boolean=non-interactive "${OPTS[@]}"; then
        terraform apply -auto-approve -var-file=<(get_vars_file_content)
    else
        terraform apply -var-file=<(get_vars_file_content)
    fi
}

function output {
    terraform output -json | jq 'map_values(.value)'
}

function destroy {
    state_tag="$(get_brick_sanitized_name "$(get_brick_name "$(pwd)")")"
    if get_arg --boolean=non-interactive "${OPTS[@]}"; then
        terraform destroy -auto-approve -var-file=<(get_vars_file_content)
    else
        terraform destroy -var-file=<(get_vars_file_content)
    fi
}

function validate {
    terraform validate
}

function fmt {
    terraform fmt
}

function show_dependencies {
    if env="$(get_env)"; then
        provider_dependencies="infra-grounds/1-init/1-$(get_env).sh"
    fi
    
    comment_dependencies="$(default_show_dependencies)"
    
    terraform_remote_state_config="$(cat "$brick_path"/*.tf |
        sed -n '/^data "terraform_remote_state" ".*" {/,/^}/ p' |
        sed -n '/  config = {/,/  }/ p')"
    local_backend_dependencies="$(echo "$terraform_remote_state_config" |
        grep "^ *path *=" | sed 's|^ *path *= *"\(.*\)".*$|\1|g')"
    gcs_backend_dependencies="$(echo "$terraform_remote_state_config" |
        grep "^ *prefix *=" | sed 's|^ *prefix *= *"\(.*\)".*$|\1|g')"
    
    echo "$(echo "$provider_dependencies" ;
        echo "$comment_dependencies" ;
        echo "$local_backend_dependencies" ;
        echo "$gcs_backend_dependencies")" |
        sed '/^$/d' | sort | uniq
}

function help {
    echo "The help of this brick haven't been overloaded so it runs the normal way."
    echo "execiac BRICK_PATH ACTION [OPTIONS]"
    echo "init: run terraform init"
    echo "plan: run terraform plan -detailed-exitcode -var-file=infra-grounds/1-init/"
    echo "apply: run terraform apply"
    echo "output: run terraform output -json"
    echo "destroy: run terraform destroy"
    echo "validate: run terraform validate"
    echo "fmt: rewrite file to pass linter"
    echo "show_dependencies: search comment #TAG:depends:BRICK_NAME"
    echo "install: install terraform in /usr/local/bin"
    echo ""
    echo "special options:"
    echo "--non-interactive: add -auto-approve opts to terraform apply"
}

