#!/bin/bash

brick_dot_yaml_path="$1"

function debug {
  label="$1"
  shift
  if [ "$(wc -l <<<"$@")" -gt 1 ]; then
    echo "DEBUG:$label<<<" >&2
    echo "$@" >&2
    echo "DEBUG:$label>>>" >&2
  else
    echo "DEBUG:$label <$@>" >&2
  fi 
}

function yaml2json {
    PYTHON_COMMAND='import sys, yaml, json; y=yaml.safe_load(sys.stdin.read()); print(json.dumps(y))'
    python3 -c "$PYTHON_COMMAND" <<<"$1"
}
function json2yaml {
  PYTHON_COMMAND='import sys, yaml, json; j=json.loads(sys.stdin.read()); print(yaml.dump(j))'
  python3 -c "$PYTHON_COMMAND"
}

function pretty_display {
  v2_json="$(yaml2json "$1")"
  echo "version: $(jq -r .version <<<"$v2_json")"
  echo "module: $(jq -r .module <<<"$v2_json")"
  echo "dependencies:"
  echo "$(jq .dependencies <<<"$v2_json" | json2yaml | sed 's/^\(.*\)$/  \1/g')"
  echo "inputs:"
  jq -c .inputs[] <<<"$v2_json" | while read line ; do
    echo "- type: $(jq -r .type <<<"$line")"
    echo "  format: $(jq -r .format <<<"$line")"
    echo "  path: $(jq -r .path <<<"$line")"
    echo "  datas:"
    json2yaml <<<"$(jq .data <<<"$line")" | sed -e 's/^\(.*\)$/    \1/g'
  done | sed '/^ *$/d'
}

v1="$(yaml2json "$(cat "$brick_dot_yaml_path")")"
v2="$v1"

input_length="$(jq '.input | length' <<<"$v2")"
for i in $(seq 0 1 $(($input_length-1))) ; do
    data="$(jq ".input[$i].data" <<<"$v2")"
    v2current_deps_yml="$(json2yaml <<<"$data" | sed 's|- name: \(.*\)$|\1:|g')"
    v2data="$(yaml2json "$(json2yaml <<<"$data" | sed -e '/^- from: .*$/d' -e 's|^ *name:|- |g')")"
    v2="$(jq ".input[$i].data=$v2data" <<<"$v2")"
done

v2="$(json2yaml <<<"$v2")"

# final v2 but the attribs are in alphabetical order and inputs.datas is named inputs.data
v2="$(echo "$v2_1";
  echo "version: 1.0.0"
  grep "^module: " <<<"$v2"
  echo "dependencies:" ; 
  jq -c .input[].data[] <<<"$v1" | while read line ; do 
    echo "  $(jq -r .name <<<"$line"):" ; 
    echo "    from: $(jq -r .from <<<"$line" | sed 's|^\([^:]*\):\($.*\)$|\1:output:\2|g')" ;
  done
  echo "inputs:" ; 
  grep -vE "^(version: |module: |input:)" <<<"$v2")"

pretty_display "$v2" > "$brick_dot_yaml_path"


# exception:
# tests/repos/infra-grounds/1-init/1-create_accounts/brick.yml
# tests/repos/infra-users/1-production/1-bastion/brick.yml
# tests/repos/infra-users/1-staging/1-bastion/brick.yml
# tests/repos/infra-users/1-monitoring/1-bastion/brick.yml
# tests/repos/infra-users/0-users_and_groups/brick.yml
# find tests/repos/ -name 'brick.yml' | grep -Ev "(/infra-grounds/1-init/1-create_accounts/|/infra-users/.*/1-bastion/|/infra-users/0-users_and_groups/)"

# v1='version: 0.0.1
# module: terraform
# input:
#   - type: env_vars
#     format: env
#     path: ""
#     data:
#       - name: TF_VAR_project_id
#         from: infra-ground/init/create_accounts:$.projects.staging.project_id
#       - name: TF_VAR_env_name
#         from: infra-ground/init/create_accounts:$.projects.staging.env
#   - type: file
#     format: json
#     path: from_exeiac.auto.tfvars.json
#     data:
#       - name: network_id
#         from: infra-ground/envs/staging/network:$.network.network_id
#       - name: network_ip_range
#         from: infra-ground/envs/staging/network:$.network.ip_range
#       - name: private_domain_name
#         from: infra-ground/envs/staging/network:$.domain_name.private
#       - name: internal_domain_name
#         from: infra-ground/envs/staging/network:$.domain_name.internal
# '

# v2='version: 0.0.1
# module: ansible
# dependencies:
#   TF_VAR_project_id:
#     from: infra-ground/init/create_accounts:$.projects.staging.project_id
#   TF_VAR_env_name:
#     from: infra-ground/init/create_accounts:$.projects.staging.env
#   network_id:
#     from: infra-ground/envs/staging/network:$.network.network_id
#   network_ip_range:
#     from: infra-ground/envs/staging/network:$.network.ip_range
#   private_domain_name:
#     from: infra-ground/envs/staging/network:$.domain_name.private
#   internal_domain_name:
#     from: infra-ground/envs/staging/network:$.domain_name.internal
# inputs:
#   - type: env_vars
#     format: env
#     path: ""
#     data:
#       - TF_VAR_project_id
#       - TF_VAR_env_name
#   - type: file
#     format: json
#     path: from_exeiac.auto.tfvars.json
#     data:
#       - name: network_id
#       - name: network_ip_range
#       - name: private_domain_name
#       - name: internal_domain_name
# '