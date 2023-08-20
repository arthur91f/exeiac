# Exeiac

## Description
`exeiac` is a tool that enables infrastructure folks to handle several
different provisionning IaC (Infrastructure as Code) tools, under a
single CLI, and helps solving some recurrent paintpoints with IaC.

It follows the brick convention, which describes an infrastructure as
a set of bricks. A brick is a piece of infrastructure, it
simultaneously is the actual infrastructure element, the code that
describes it and a piece code that allows its execution; be it a
terraform state, an ansible playbook or a helm chart for instance.

This project was born from the following needs:
- solve dependencies issues
- increase transparency as to how a piece of infrastructure should be
    deployed. No matter the provisionning tool you use, if there's a
    module for it, `exeiac` will handle it
- allow for a clean way to interact with only a part of your
  infrastructure in a safe way, without breaking the dependency tree,
  even if your infrastructure management tool doesn't provide that
  feature

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

## Get started

### Installation

Clone the git repository and build:
``` bash
$ git clone github.com/arthur91f/exeiac/src/exeiac
$ cd exeiac
$ go install src/exeiac
```

There is no release process yet, but on Go version 1.16 or later you can:
``` bash
# Install at tree head:
$ go install github.com/arthur91f/exeiac/src/exeiac/src/exeiac@main
```

- get the exeiac binary
- have an infra code that follow some conventions (see below for more details)
  - each brick is a directory prefixed by a number to make the apply order
    transparent
  - each elementary brick should have a brick.yml to define how it will be
    executed and the input it needed from dependencies
  - each elementary brick should reference in brick.yml an executable or module
    to execute itself. (a module is simply an executable that is not in the
    brick directory and that can be called by many bricks)
  - Read [howto_write_brick](./docs/howto_write_brick.md)
- create a conf file in /etc/exeiac/exeiac.yml or $HOME/.config/exeiac.yml
  ```yaml
  modules_path:
    terraform: $HOME/git-repos/exeiac-modules/terraform
    ansible: $HOME/git-repos/exeiac-modules/terraform
  room_paths_list:
    - $HOME/git-repos/infra-ground
    - $HOME/git-repos/applications
    - $HOME/git-repos/users
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

Read [howto_write_module](./docs/howto_write_module.md)
