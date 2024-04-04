# Writing Your Module

## Generality

A module is an executable that will be called by exeIaC to run an action. The module will implement an interface. Some of the methods of this interface have to be implemented, while others can be omitted. You can also overload this interface with additional actions. These actions will be described more precisely later, but here are some examples:
- **describe_module_for_exeiac**: mandatory, used by exeIaC to get actions implemented by the module.
- **lay**: mandatory, to deploy or correct drift on a brick.
- **remove**: mandatory, rollback a lay. This action is less mandatory than the previous two because you can manage your infra by removing manually.
- **output**: display output information of a brick (mandatory if you want some brick depends on a brick using this module).
- **plan**: check if a lay is needed and/or what the lay will do.
- **help**: display specific help for the module or more specific for the brick.
- **init**: install or check all prerequisites for the module.
- **clean**: remove all files created by the module.
- **validate_code**: validate the syntax of the brick's code.
- ... it can be overloaded

**Note: module actions aren't exeiac commands!**
You may recognize some exeIaC actions, but they are not exactly the same. For example, the module lay can be just a _terraform apply_. The exeIaC _exeiac lay_ will of course execute the module lay, but before it will search all the bricks it needs to output to get the input of the lay. Then it will output the current brick and register this output and after the module lay, it will output again the current brick to compare the output before and after and determine if it has changed. But happily it's exeIaC that does all of that; you just have to code the module lay and output and specify the inputs needed.

So you haven't to implement:
- show
- get-depends
- smart-lay

## Think to specify your module in your exeIaC and bricks conf

Find a name for your module. It will be used in:
- brick.yml files: at field .module
- exeIaC conf: at field modules[].name
  and specify the full path of your module

## Module's Inputs Conventions

- **Brick Path**: before calling module exeIaC will change directory to the brick path.
- **Action**: module takes one positional argument that is the *action* name.
- **Options**: module can take some other options as follows:
    - --non-interactive (is passed to module by exeiac interpreting -i, -interactive, -I or --non-interactive options).
    - other options of your invention can be passed with -o or --other-options.
- **exeIaC Env Variables**: exeIaC defines some env variables that can be used by modules:
    - EXEIAC_BRICK_PATH
    - EXEIAC_BRICK_NAME
    - EXEIAC_ROOM_PATH
	- EXEIAC_ROOM_NAME
	- EXEIAC_MODULE_PATH
	- EXEIAC_MODULE_NAME
- **Brick Inputs**: brick can define some inputs. It can be a file or env vars (look at howto_write_brick_yaml.md).

## Module's Outputs Conventions

The different actions have different behavior. For some, the stdout will be displayed; for others, it will be registered inside. For some, the status code will be interpreted; for some others not.

### Stdout

Except for *describe_module_for_exeiac* and *output*, where the stdout will be consumed by exeIaC (and so need to follow some convention), all other actions will just display the stdout without modifying it.

For output and describe_module_for_exeiac, the expected format is json. For output, no other convention is needed. For describe_module_for_exeiac, check the proper part.

### Stderr

The stderr will always be displayed without being consumed.

### Status Code

By default, the status code understanding will be that:
- 0: when the action succeeds.
- 1-255: when a problem occurs (it will lead exeIaC to stop runs).

But you can define some events that won't be considered as a fail for an action by setting displayed with describe_module_for_exeiac. For example, 0 corresponds to no drift and 2 corresponds to there is/was a drift.

### describe_module_for_exeiac stdout

It displays a dictionary of implemented actions: Each action is a dictionary with those fields:
- behaviour: that can be omitted nowadays because exeIaC only implements one valid default behavior by action. But it may change in the future.
- status_code_fail: default is "1-255" but you can change it.
- events: dictionary that contains those fields:
    - type: 
        - *status_code*: boolean set to true if it corresponds to the field status_code.
        - *file*: string represent the content of the file.
        - *json*: the content of the file can represent whatever you want as long as it is in json format.
        - *yaml*: the content of the file can represent whatever you want as long as it is in yaml format.
    - status_code: only for type=*status_code* that is a number sequence as "2,5-10".
    - path: only for type=*file*, *json*, *yaml*. It can be a full path or a relative path from the brick path. It can also be a classic file or a named pipe.

    For plan exeiac can interpret some special event to know if there is a drift or not:
    - for type *status_code*:
        - exeiac_plan_no_drift
        - exeiac_plan_drift
        - exeiac_plan_unknown
    - for type *file*, *json* or *yaml*: the event name should be exeiac_plan and the content should be in equal to:
        - no_drift
        - drift
        - unknown
    A success will be considered as an exeiac_plan_unknown except if a previous event is caught.


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
