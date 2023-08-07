# Writting your module

## Generality

A module is an executable that will be called by exeIaC to run an action.
The module will implement an interface. Some of the method of this interface
have to be implemented some other can be omitted. You can also overload this
interface with some other actions.
These actions will be describe more precisely later but here some :
- **describe_module_for_exeiac**: mandatory, used by exeIaC to get action
  implemented by the module.
- **lay**: mandatory, to deploy or correct drift on a brick
- **remove**: mandatory, rollback a lay. Is less mandatory than the 2 previous
  action, because you can manage your infra by removing by hand.
- **output**: display output information of a brick
    (mandatory if you want some brick depends of brick using this module)
- **plan**: check if a lay is needed and/or what the lay will do
- **help**: display a specific help for the module or more specific for the
    brick
- **init**: install or check all prerequisite for the module
- **clean**: remove all files created by the module
- **validate_code**: validate teh syntax of the brick's code
- ... it can be overloaded

**Note: module actions aren't exeiac commands !**
You may recognize some exeIaC actions but they are not exactly the 
same. For example the module lay can be just a _terraform apply_. The exeIaC
_exeiac lay_ will of course execute the module lay but before it will search 
all the bricks it needs to output to get the input of the lay. Then it will
output the current brick and register this output and after the module lay
it will output again the current brick to compare the output before and after 
and say if it has changed. But happily it's exeIaC that do all of that, you
just have to code the module lay and output and specify the inputs needed.


## Module's inputs conventions

- **Brick path**: before call module exeIaC will change directory to brick path
- **Action**: module takes one positionnal argument that is the *action* name
- **Options**: module can takes some other options as following
    - --non-interactive (is passed to module by exeiac interpretating -i, 
        -interactive, -I or --non-interactive options)
    - other option of your invention can be passed with -o or --other-options
- **exeIaC env variables**: exeIaC define some env variables that can be used by modules
    - EXEIAC_BRICK_PATH
    - EXEIAC_BRICK_NAME
    - EXEIAC_ROOM_PATH
	- EXEIAC_ROOM_NAME
	- EXEIAC_MODULE_PATH
	- EXEIAC_MODULE_NAME
- **Brick inputs**: brick can define some inputs. It can be a file or env vars 
    (look at howto_write_brick_yaml.md)


## Module's outputs conventions

The different actions have different behaviour. For some the stdout will be 
displayed, for other it will be registered inside. For some the status code 
will be interpretated, for some other not.

### Stdout

Except for *describe_module_for_exeiac* and *output* where the stdout will be 
consumed by exeIaC (and so need to follow some convention), all other actions 
will just display the stdout without modify it.

For output and describe_module_for_exeiac, the expected format is json.
For output no other convention is needed, for describe_module_for_exeiac, check
the proper part

### Stderr

The stderr will always be displayed without beeing consumed

### Status code

By default the satus code understanding will be that
- 0: when action succeed
- 1-255: when a problem occurs (it will leads exeIaC to stop runs)

But you can define some events that won't be consider as a fail for an action
by setting diplayed with describe_module_for_exeiac.
For example 0 correspond to no drift and 2 correspond there is/was a drift.

### describe_module_for_exeiac stdout

It displays a dictionnary of implemented actions:
Each action is a dictionnary with thoose fields:
- behaviour: that can be ommitted nowaday because exeIaC only implement one 
    valid default behaviour by action. But in may change in future.
- status_code_fail: default is "1-255" but you can change it
- events: dictionnary that contain thoose fields
    - type: 
        - *status_code*: boolean set to true if it correspond to field 
            *status_code*
        - *file*: string represent the content of the file
        - *json*: the content of the file can represent whatever you want as
            long as it is in json format
        - *yaml*: the content of the file can represent whatever you want as
            long as it is in yaml format
    - status_code: only for type=*status_code* that is a number sequence as 
        "2,5-10"
    - path: only for type=*file*, *json*, *yaml*. It can be a full path or a relative 
        path from the brick path. It can also be a classic file or a named pipe.

    For plan exeiac can interpret some special event to know if there is a drift 
    or not:
    - for type *status_code*:
        - exeiac_plan_no_drift
        - exeiac_plan_drift
        - exeiac_plan_unkown
    - for type *file*, *json* or *yaml*: the event name should be exeiac_plan 
        and the content should be in equal to:
        - no_drift
        - drift
        - unknown
    A success will be considered as an exeiac_plan_unkown except if a previous event
    is catched.

Example 1:
```json
{
    "lay": {
        "status_code_fail": "1,4-255",
        "events": {
            "drift_corrected_without_recreation": { 
                "type": "status_code",
                "status_code": 2
            },
            "resources_recreated": { "type": "status_code", "status_code": 3 }
        }
    },
    "remove": {},
    "output": {},
    "plan": {
        "status_code_fail": "1,3-255",
        "events": {
            "exeiac_plan_no_drift": { "type": "status_code", "status_code": 0 },
            "exeiac_plan_drift": { "type": "status_code", "status_code": 2 }
        }
    },
    "help": {},
}
```

Example 2:
```json
{
    "lay": {
        "events": {
            "modified_resources_list": { 
                "type": "json", 
                "path": "./.modified_resources.json"
            },
            "recreated_resources_list": { 
                "type": "json", 
                "path": "./.recreated_resources.json"
            }
        }
    },
    "remove": {},
    "output": {},
    "plan": {
        "status_code_fail": "1-255",
        "events": {
            "exeiac_plan": { "type": "file", "path": "./.exeiac_plan" },
        }
    },
    "clean": {},
    "help": {}
}
```
