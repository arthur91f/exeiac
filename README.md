# exeiac

execute infra code smartly by :
- trying to solve dependency issues
- being transparent on how to deploy a part of your infra (no matter if you
  use terraform or ansible or... you will use the same interface)
- permit proper shortcut that don't break dependency tree even if it's not
  easily managed by your tool

exeiac use the brick metaphor. It describes infra as a set of bricks. A brick
is a part of infra (the real infra element, the code that describes it, the code
that permits to execute it). It can be a terraform state, an ansible playbook,
a helm chart...

## DOCUMENTATION

- philosophy: theoretical approach useful to write your infra code and understand
  exeiac best practices
  - common infra code problems we try to solve with exeiac
  - define infra as a set of bricks
  - define the brick concept
  - explain how bricks depends of each other
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
  ```yaml
  modules_path:
    terraform: $HOME/git-repos/exeiac-modules/terraform
    ansible: $HOME/git-repos/exeiac-modules/terraform
  room_paths_list:
    - $HOME/git-repos/infra-ground
    - $HOME/git-repos/applications
    - $HOME/git-repos/users
  ```

### Compile

Once cloned and in the repo directory, you'll be able to compile the project by
running the following command:
```bash
go build src/exeiac
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
- deploy a brick and recursively deploy all bricks that depends on an output
  that have changed. Note that here we have used the brickname and not the path
  ```bash
  exeiac lay infra-core/staging/ssh_bastion --bricks-specifier=selected+needed_dependents
  ```
- destroy a higher level brick. It will destroy all elementary bricks
  contained in the higher level bricks in the right order.
  ```bash
  exeiac remove infra-core/staging
  ```
- get more help
  ```bash
  exeiac help
  ```

### Create a module or an executable

#### Synopsis

module-executable ACTION PATH [OPTIONS]

#### Actions to implement

Your module can implement every standard exeiac action method for a brick,
but exeiac will have a default action if you haven't implemented the specified
action. You can also implement other personal action.

- action to implement. If not it will display an error or don't follow the
  exeiac logic when you call it:
  - show_implemented_actions: will display an error for any action if it's not
    implemented
  - plan:
    if it's not implemented: will always assume that there is a drift
    exit_code:
    - success with no drift: 0
    - success with drift: 2
    - success when module can't decide between 0 and 2: 3
    - fail: 1,4-255
  - lay:
    if it's not implemented: will always assume that the deploy has failed
  - remove:
    if it's not implemented: will always assume that the destroy has failed
  - output:
    if it's not implemented: the ouput will be null: {}
    format: json
- classic action you can overload, if not, they will simply do nothing:
  - init: install or check tools dependencies (terraform binary, provider,
    python, ansible library...)
  - validate_code
  - help
  - personnal_action
- action to not implement, because they are already in the exeiac logic:
  - show --format input: use brick.yml
  - show --format direct-previous: display the list of bricks whom output are
    needed for execute the brick
  - show --format direct-next: display the list of bricks that needs the output
    of the specified bricks
  - show --format linked-previous: display the list of bricks that needs to be 
    layed to lay your brick (as their outputs are needed directly or not)
  - show --format linked-next: display the list of bricks that can be impacted 
    recursively by a change in your brick
  - show --children: display elementary bricks 
  - ... others exeiac command in general
