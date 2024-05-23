# A new approach for these tests

Here we will try to end to end tests with a module that answer as the brick decide


## Module
We can develop a debug module that just write debug.txt
```
#######
# ARG #
...

#######
# ENV #
<env var that correspond to regex of option --debug-env>

##############
# INPUT FILE #
...

############
# <ACTION> #
...

#######
# ENV #
...
```

## Brick

A test brick will be very simple and descriptive. The output will not be 
function of input. The output 
- brick.yml file (comment explaining the behaviour if the brick is supposed to be broken or not)
- state.exeiactest.json correspond to output
- code.exeiactest.json correspond to the desired output
- special_behaviours.exeiactest.json to specify a special action behaviour.
```json
{
    "lay": {
        "stdout": "TEST: begin lay",
        "stderr": "TEST: an error is emulated",
        "final_state": {},
        "return_code": 1,
        "events": {},
        "run": [ "stdout", "stderr", "default", "final_state", "return_code", "events" ]
    }
}
```

## Infra

### Tree
repos
├── apps
│   ├── 1-database
│   └── 2-myapp
│       ├── 1-instance
│       └── 2-configuration
├── infra
│   ├── 1-base
│   ├── 2-network
│   └── 3-bastion
│       ├── 1-instance
│       └── 2-configuration
└── modules
    └── special_behaviour_examples

### Inputs
- infra/base
- infra/network
  - infra/base:output:$.cloud_provider.project_id
  - infra/base:output:$.cloud_provider.admin_token
- infra/bastion/instance
  - infra/base:output:$.cloud_provider.project_id
  - infra/base:output:$.cloud_provider.admin_token
  - infra/network:output:$.network_id
  - infra/network:output:$.ip_range
- infra/bastion/configuration
  - infra/bastion/instance:output:$.ip
- apps/database
  - infra/base:output:$.cloud_provider.project_id
  - infra/base:output:$.cloud_provider.admin_token
  - infra/network:output:$.network_id
  - infra/network:output:$.ip_range
- apps/myapp/instance
  - infra/base:output:$.cloud_provider.project_id
  - infra/base:output:$.cloud_provider.admin_token
  - infra/network:output:$.network_id
  - infra/network:output:$.ip_range
- apps/myapp/configuration
  - infra/bastion/instance:output:$.ip
  - infra/bastion/configuration:output:$.users.myuser
  - apps/myapp/instance:output:$.ip
  - apps/myapp/instance:output:$.ram
  - apps/database:output:$.instances.master.address
  - apps/database:output:$.users.myapp

## Create a new test

### Write an item in a yaml file in tests/tests:
**title**: a string to describe the test
**cmd**: the command to execute for the test
**status**: the status code that should return the command
**stdout**: the strings that the command should display

