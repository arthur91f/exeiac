```markdown
# Write Your Brick

## Generality

A brick is an element of your infra and its IaC representation. It can be a Terraform state, an Ansible playbook, a Helm chart, a combination of all those things, or even other elements depending on your brick's module.

A brick is represented by a directory.

But a brick is not only the code describing the real infra element. The elements describing a brick are:
- code: content of the directory
- outputs: the state of the brick, the values taken by things declared in codes as generated passwords, given public IP addresses, or just some values that can be used by other bricks
- inputs: describe the context of the brick, it is composed of parts of outputs of other bricks
- module events: output of the module action, they are describing what has been done not the state of the brick (that is described by outputs)
- triggers: input change or other brick module's event that will trigger a lay or another action on the brick

**Note: Elementary Bricks vs Higher-Order Bricks:** Every brick is represented by a directory. Some bricks contain other bricks and some others just contain a brick.yml file that describes its module, inputs, triggers...

## Directory Convention

You can identify bricks from their path:
- A brick must be in a "room" (listed in [exeiac configuration file](./howto_write_configuration_file.md)) or be a room.
- A brick directory name must begin with a priority number that describes the building sequence of bricks inside the same room or super brick.

**Note: Why the Priority Number?** Actually, it's redundant with dependency information, and exeIaC should be able to build its valid laying sequence. But we wanted to have a convention that lets humans easily partially guess some layer and dependency tree. It is also useful to identify easily what directory is a brick. Based on this, you can add a documentation directory or a directory with maintenance scripts like dump and restore database that will be ignored by exeiac.

## Writing brick.yml File

In each elementary brick directory, you have to write a brick.yml file. It needs to contain these fields:
- **version**: the version of the file. At the moment, we are in 1.0.0
- **module**: the module name or a relative path to the module beginning by ./
- **dependencies**: a dictionary of input values needed to execute the brick.
  - **from**: the format is brickname:source:jsonpath
    - _brickname_: the brick name of the dependencies
    - _source_: can be "output" or "event"
    - _jsonpath_: the JSON path of the value, for example: *$.network.ip_range*
  - **needed_for**: list of actions that need that dependency. The default is set in exeiac configuration files
  - **not_needed_for**: list of actions that don't need that dependency. The default is set in exeiac configuration files
  - **triggered_actions**: what action will be triggered by that dependency. Default is lay
  - **trigger_type**: default is classic. Can be:
    - _classic_: if the value changes, it will trigger the action in *trigger_action*. To use if the input is also a spec of your brick
    - _weak_: if the value changes, it won't trigger anything. To use if the input is needed for execution but is not a spec of your brick. For example, the hop SSH server credentials used to connect to your VM.
    - It can also be a special trigger type usually used with events. It will be a function as:
      - is_true()
      - is_false()
      - is_not_null()
      - is_null()
      - is_equal(string)
      - is_different(string)
      - contains(string)
      - not_contains(string)
      - match(string)
      - not_match(string)
- **inputs**: how to present dependencies value to the module. see below
  - **type**: can be env_vars or file
  - **format**: can be env or json
  - **path**: the relative path if it's type file else left to default ""
  - **datas**: the list of the dependencies keys
```
Ce texte peut être copié et collé dans un fichier Markdown pour être utilisé.