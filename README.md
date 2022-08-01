# exeIaC

Execute infra code according to a bricks dependencies metaphore

## METAPHORE EXPLANATION

### Brief explanation

We understand IaC as a construction composed with set of bricks. A brick is a step for building the IaC, it can be a terraform directory, an ansible playbook, a kubernetes yaml file...
Most of those bricks depends of others. Actually you can't put bricks of the second floor before building the first floor. And same you can't add a user to a database before poping the database instance on the cloud provider.


### Lot of infra codes problem

In lot of company the infra code division in bricks and the right to apply them were often hardly linked to technology and the dependencies between bricks are not clear.
In a bit complexe infra code it brings :
- mistakes of infra team when applying a change (because they don't understand dependencies) and they often guess how it should be apply.
- a big documentation not always up to date for each detail specificities (as the dependency tree is not clear)
- shortcuts where you copy paste a generated value from an other brick (because convention exchange info between bricks doesn't exist)
- some infra code omission because it can't be done (or it's to complicated) by terraform
- interrogation without answer when you discover an instance on your webUI cloud provider and you don't kow where it has been coded.
- to avoid drift and CD complexity, all infra is reapplied regularly. It makes development less fluent and sometimes the CD failed due to a teamate modification that you have to correct before see your change applied on your brick (a bit annoying when you try to fix a production incident)

ExeIaC try to answer to those problems.

### ExeIaC philosophy

The main philosophy of exeIaC is to decouple the brick technology complexity to the way bricks interact with each other. It try to respond to the dependency inversion principle of SOLID.
How to decouple ? More than a script that manage your infra code, exeIaC is a convention and each brick have to respect this convention. In fact the brick contract is to be able to be applied, display output and display dependencies bricks... (look forward to see the exact convention)

What brings this decoupling ? Each brick can be applied in a different way with different technologies. We don't recommend you to use too many different way to apply your bricks  (ExeIaC implement modules to factorize application code for brick of the same type), but it permit you shortcut.

For example, if you need to open an ssh tunnel or switch to another vpn before running your terraform you can do that easily and you won't need to document that tricky specificities because it will be written clearly in the code.

In the same way if you have to do something by hand you can write it in your brick in place of execute a terraform apply it will get information from other bricks, display some instructions to do by hand and then display outputs to the dependencies. That type of shortcut has to be considered ugly. But they're still better than write nothing. It won't be invisibilize, it won't break any dependencies, it will warn anyone trying to applied the brick. Moreover, in most case it can be scripted and you won't need to display instruction to do by hand.

This decoupling will permit to reduce the documenation size and apply the rule "no details or specificities written in the documentation".
Why this rule ? Because documentation is hard to maintain up to date and exhaustive. And a short general up to date documentation will always be better than an attempt of exhaustivity where you don't know how to search something, where you're not sure that it exist and where you not sure that it is up to date if it exist.
How this decoupling will allow to reduce documenation ? Because you will document your brick inside your brick directory (where the code is written). The detail documentation will be easy to found as your tree structure should be intuitive. And when update the brick it is to update a README.md file if needed and to block the PR until it has been done. Moreover when you debug, the documenation will be in the same directory !

As everything is coded (thanks to allow easy shortcut) and the dependencies between bricks is clear the mistakes by applying something are reduce and the impact of a fail is well known.

You can also avoid to discover cloud provider resources that you can't link to a code by following a best practice: the path of your brick has to be tagged on the resources. 


### IaC vocabulary

- brick: a part of the infra that can be managed unitary. It can be the configuration of a monitoring user on a database or the creation a an entire cluster. A brick as a part of a wall is on top of other bricks that can be considered as dependencies and some brick can be on top of it that are the dependents. As we are in infra as code a brick is describe by code and the code can be execute. By convention we can execute some task on a brick :
  - init (install tools, download library...)
  - plan (check what code will change on existing)
  - apply (apply the code)
  - output (return some value use by dependents bricks)
  - destroy (revert the apply)
  - validate (check that code is correct)
  - fmt (linter, able to rewrite files)
  - show_dependencies (show dependencies bricks)
- super brick: a super brick is a brick that contains other bricks it can be applied by apply sequentially all subricks it contains by following the dependency tree recursively.
- elementary brick: a brick that is not a super brick
- room: a super brick that is subjective cut of the infra that make sense. Generally a team have a specific ownership on a room and it can be assimilate to a domain (in DDD overview). But in a simple code overview lets say that a room is a git repository.
- other brick regroupment :
  - floor you can talk about floor the first floor is the room infra transverse to all domains. on the second floor you have the room client, backoffice, AI... it's the floor used by developers and feature teams. The third floor is transverse to all domains and manage users access on the infra.
the floor abstraction should simplify the dependencies abstraction.
  - pillar
- furniture: a set of scripts to use in daily usage but that isn't use for building your infra. It can be a script to dump a databae. It is a good idea to put those sort of script on your infra as code so when you have to debug you have all your tools and documentation in the same place but obviously they are completely optionnal and they have nothing to do with dependency tree.
brick path: the full path of a brick on your computer
brick name: the name of the repo and the relative path from nome of repo to the brick
brick name sanitized: if you use brick name as tag on your cloud resources you may have to change slash 
 to underscore or dash. it's up to you to define a the sanitized regex

## ExeIaC

### Short introduction to bash

- there is no type in bash, everything is a string (it is little more complicated but if you begin just remember that)
- all commands, function, scripts... return an int (0-255). 0 means the execution is ok, the other numbers are error code. So all commands can be tested as a boolean.
- If you want your command return a string just make your command display it on standard output by using echo command
- Usually we don't manage array in bash we just use string with space for i in a b c d ; do ...
- >&2 redirct output on error output
- <<<"mystring" can be put inplace of an arg that expect a filepath
- result="$(my_command)" ; err_code=$? store what is display by my_command in variable result and the error code returning by my_command in err_code
- grep "^foo" path/bar  <=>  grep "^foo." <<<"$(cat path/bar)"  <=>  cat path/bar | grep "^foo"  <=>  grep "^foo." <(cat path/bar)
- $1 $2 ... ${10} ... are positional arguments "$@" represent the list of positionnal argument
- ${a} <=> $a

- Take care to " especially if you want conserve line break or argument number.
- take care to space in if [ condition ] statement.

### exeIaC tool code vocabulary

**module**: a module is a set of bash function corresponding to bricks actions (all actions haven't to be present as is used overwriting). each module should implement an additionnal function : is_brick_using_this_module
**action**: a brick action as:
- init (install tools, download library...)
- plan (check what the brick code will change on existing resources)
- apply (apply the code)
- output (return some values used by dependents bricks)
- destroy (revert the apply)
- validate (check that code is correct)
- fmt (linter, able to rewrite files)
- show_dependencies (show dependencies bricks)

**dependencies**: is the need that an other brick is already applied or some previous brick output.

## CONVENTION & BEST PRACTICES

### Convention
All modules should implement the same action describe up. The simple way to do that is to heritate from the default module (source $modules_path/default.sh) so if you haven't implement fmt, you won't have bugs.
All actions should be idempotent.
All dependencies have to be visible. The default module permit to track all dependencies commented like that #EXEIAC:depends:BRICK_NAME or //EXEIAC:depends:BRICK_NAME
Then choose one unique convention beside for all your infra code:
- Apply action is always: plan, ask confirmation, apply. All interactive action should understand --exeiac-opts=non-interactive to skip plan and ask confirmation.
- No action should be interactive (the code will be simplier if you haven't to implement non-interactive option)
All brick are directory or file with name begins by a figure and a dash (2-database). It permit to identify brick and brick content prioritize the apply of different brick in the same directory and to add directory and file that are not re

### Best practices 

1. Don't multiply modules that use different technology. Try to keep the applying infra code easy and intuitive. You should know how to apply a brick without reading its module ( Because it is always the sames ). We tend to say that an average infra collaborator should not masterize more than 3 modules and that masterize one complexe module is too much for most of developers.
However you can use heritage to get multiple modules that use a different terraform version. If they keep the same logic and it can be clearer than one unique complexe module that trying to manage different binary version.
Don't hesitate to display warning messages or asking confirmation when you implement a tricky shortcut.
2. Tag your cloud resources with the brick_name or sanitized_brick_name so when you will see something in cloud it will be very easy to find it in code.
3. All outputs have to present the same format like json. The more you will define convention about outputs the more pass on change on dependents will be easy after modifying a brick.

## DISCUSSION

### Why bash

For the first version, the aim is to execute bash instrcution write in console so use bash should be the more simple way to do that to conserve coloration output of tools etc.
Moreover, even if bash isn't the easiest language to develop, it's not a such big deal for me, and it will be more simple to debug module for others as you just have to copy paste code on terminal.
Modules can call bash exeiac functions if needed. But better try to call them with exeiac. 

### Performance aim:

- exeIaC runtime to execute one bricks should be instant and very cheap on resources
- exeIaC runtime to plan its execute plan should be less than 3 seconds
- exeIaC should not open to much bash in the same time ithout parallelism
- enable parallelism

