# Exeiac

## Why use exeIaC
`exeiac` is a tool that enables infrastructure folks to handle several 
different provisioning IaC (Infrastructure as Code) tools under a single CLI 
and helps solve some recurrent pain points with IaC.

It follows the brick convention, which describes an infrastructure as a set of 
bricks. A brick is a piece of infrastructure; it simultaneously is the actual 
infrastructure element, the code that describes it, and a piece of code that 
allows its execution; be it a terraform state, an ansible playbook, or a helm 
chart, for instance.

This project was born from the following needs:
- solve dependencies issues
- increase transparency as to how a piece of infrastructure should be deployed. 
  No matter the provisioning tool you use, if there's a module for it, `exeiac` 
  will handle it
- allow for a clean way to interact with only a part of your infrastructure in 
  a safe way, without breaking the dependency tree, even if your infrastructure 
  management tool doesn't provide that feature

## What is a brick
First, exeIaC deals with infra bricks. What is an infra brick? Basically, it's 
a directory that contains some infra code. It can be a terraform directory to 
deploy a VM, an ansible playbook to configure a host, a helm chart, or simply a 
template that describes some instructions to do manually.

This summary should be enough for you to be able to read the following 
paragraphs and have a basic understanding of exeIaC. But if you want to get a 
better grasp of the brick concept itself and its genius simplicity, you can 
find a deeper insight on 
[this documentation page](./docs/brick_concept_and_dependencies).

## Get started

- **1. Get exeiac binary**
  ```bash
  go install github.com/arthur91f/exeiac/src/exeiac/src/exeiac@main
  ```

- **2. Write your modules** in whatever language you want. A module can be seen 
  as a makefile to deploy your brick. Basically, it's a shell script that 
  follows some conventions described here: 
  [How to write module](./docs/howto_write_module.md). You have to implement 
  three commands:
    - *describe_module_for_exeiac* (that displays a json)
    - *lay* to deploy your brick
    - *output* to display some specs of your brick such as IP address, login...
    - ... you can implement other commands like plan, remove, lint, help...

- **3. Put a YAML file in each IaC directory** to describe your bricks as here 
  [How to write brick](./docs/howto_write_brick.md). It will let:
  - exeiac identify your IaC directory as an infra brick
  - associate your brick to its module
  - exeiac understand your brick's dependencies and how to present them to the 
    brick

- **4. Write your exeiac conf file** in your home or in /etc to let exeiac 
  binary find your module and your infra code. 
  [How to write config file](./docs/howto_write_configuration_file.md)

- **5. Enjoy exeiac**
  Here are some examples of basic commands you can execute:
  - Display output of a brick
    ```bash
    exeiac output infra-ground/envs/staging/network
    ```
  - Plan all sub-bricks of infra-ground/envs/staging
    ```bash
    exeiac output infra-ground/envs/staging
    ```
  - Display all bricks that should be re-deployed after the change of the 
    brick's network output .network.ip_range
    ```bash
    exeiac get-depends infra-ground/envs/staging/network -j $.network.ip_range
    ```
  - Display all bricks that can be impacted by the re-deploy of the brick network
    ```bash
    exeiac show infra-ground/envs/staging/network --bricks-specifiers linked_next --format name
    ```
  - Deploy/re-deploy a drift in brick staging/network and all other bricks 
    impacted by that drift recursively.
    ```bash
    exeiac smart-lay infra-ground/envs/staging/network --non-interactive
    ```

## Useful links

- **local**
  - [How to join us](./docs/to_write.md)
- **external**
  - [Download exeiac](https://download-exeiac.91f.ovh)
  - [exeIaC presentation](https://drive.google.com/blabla)
