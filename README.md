# exeiac

execute infra code smartly by :
- trying to solve dependency issues
- beeing transparent on how to deploy a part of your infra (no matter if you 
  use terraform or ansible or... you will use the same interface)
- permit proper shortcut that don't break dependency tree even if it's not 
  easily managed by your tool

exeiac use the brick metaphore. It describe infra as a set of bricks. A brick
is a part of infra (the real infra element, the code that describe it, the code
that permit to execute it). It can be a terraform state, an ansible playbook,
an helm chart...

## DOCUMENTATION

- philosophy: theorical approach useful to write your infra code and understand 
  exeiac best practices
  - common infra code problems we try to solve with exeiac
  - define infra as a set of bricks
  - define the brick concept
  - explane how bricks depends of each other
  - vocabulary
- development: contain specs and schema to understand how it is coded
- user: contain all you need to use the tool and create an infra code that 
  respect the convention and best practices
- examples: examples of simple infra code and module

## HOW TO

For using exeiac you will have to:
- get the exeiac binary
- have an infra code that follow some conventions (see below for more details)
  - each brick is a directory prefixed by a number to make the apply order 
    transparent
  - each elementary brick should have a brick.yml to define how it will be 
    executed and the input it needed from dependencies
  - each elementary brick should reference in brick.yml an executable or module
    to execute itself. (a module is simply an executable that is not in the 
    brick directory and that can be called by many bricks)
- create a conf file in /etc/exeiac.conf or $HOME/.exeiac.conf
  ```
  modules_path:
    - "$HOME/git-repos/exeiac-modules"
  room_paths_list:
    - $HOME/git-repos/infra-ground
    - $HOME/git-repos/applications
    - $HOME/git-repos/users
  ```

### Install

```
```

### Simple command line examples

- display a brick output
  ```json
  $ exeiac output ./infra-core/2-staging/2-ssh_bastion
  {
      "instance_id": "bastion-staging-221022",
      "private_ip": "10.11.3.2",
      "public_ip": "34.33.31.30"
  }
  ```
- deploy a brick and recursively deploy all bricks that depends of an output 
  that have changed. Note that here we have use the brickname and not the path 
  ```bash
  exeiac deploy infra-core/staging/ssh_bastion --bricks-specifier=selected+needed_dependents
  ```
- destroy an higher level brick. It will destroy all elementary bricks
  contained in the higher level bricks in the right order. 
  ```bash
  exeiac destroy infra-core/staging
  ```
- get more help
  ```bash
  exeiac help
  ```

### Create a module or an executable

#### Synopsis

module-executable ACTION PATH [OPTIONS]

#### Actions to implement

Your module can implement every standard exeiac action method for a brick. 
But exeiac will have a default action if you haven't implemented the specified 
action. You can also implement other personnal action.

- action to implement, if not it will display an error or don't follow the 
  exeiac logic when you call it:
  - show_implemented_actions: will display an error for any action if it's not 
    implemented
  - plan:
    if it's not implemented: will always assume that there is a drift
    exit_code:
    - success with no drift: 0
    - success with drift: 1
    - fail: 2-255
  - lay:
    if it's not implemented: will always assume that the deploy have failed
    exit_code:
    - success with no drift: 0
    - success with drift: 1
    - fail: 2-255
  - remove: 
    if it's not implemented: will always assume that the destroy have failed
    exit_code:
    - success with nothing to do: 0
    - success with something to do: 1
    - fail: 2-255
  - output:
    if it's not implemented: the ouput will be null
    format: json
- classic action you can overload, if not, they will simply do nothing:
  - init: install or check tools dependencies (ansible 
  - validate_code
  - help
  - personnal_action
- action to not implement, because they are in the exeiac logic:
  - show_input: use brick.yml
  - show_dependencies: use show_input
  - show_dependents: use show_dependencies
  - show_dependencies_recursively: use show_dependencies
  - show_dependents_recursively: use show_dependents
  - list_elementary_bricks: is reprensented by the directory tree
  - ... others exeiac command in general

