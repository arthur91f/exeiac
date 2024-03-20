# Write the exeiac configuration file

## Generality

The exeIaC configuration file is used to by the exeiac binary to locate your 
modules and rooms (repositories that contains your brick).

It is also used to pass some default arguments.

The configuration file is quite critical. Indeed, some information inside are 
sticky to your host some others are sticky to your infra code.
(where do you have download your infra code ?)
 (the module and room name, the action that call 
)

if you change the room's name
or the module's name, your 

## Where to write it ?

You can write the file in your home : $HOME/.config/exeiac/exeiac.yml
More generally, exeiac use the golang xdg library function xdg.SearchConfigFile
 to find the configuration file.

You can also pass the configuration file with an option: 
```bash
exeiac -c ./exeiac-conf.yml list-bricks
exeiac --configuration-file ./exeiac-conf.yml list-bricks
```
In case the documenation is not up to date try exeiac help ;)

## How to write it ?

Use the yaml format, 


## An example
```yaml
modules:
  - name: ansible
    path: /home/arthur91f/git/infra-library/exeiac-module-ansible.sh
  - name: terraform
    path: /home/arthur91f/git/infra-library/exeiac-module-terraform.py
rooms:
  - name: ground
    path: /home/arthur91f/git/infra-ground/environment
  - name: applicationA
    path: /home/arthur91f/git/applicationA
  - name: user-access
    path: /home/arthur91f/git/infra-user-access
default_arguments:
  other_options: [] # useful if you don't want to always pass the sae option as --format=name or --non-interactive

  # Will ask input of 
  default_is_input_needed: false # By default exeiac won't search the inputs of a brick it's not usefull for action init for example
  inputs_needed_for: # will search the inputs for all those action
    - plan
    - lay
    - remove
```
