# exeIaC

Execute infra code according to a bricks dependencies metaphor

---

## PHILOSOPHY

### Brief metaphore explanation

We understand IaC as a construction composed with set of bricks. A brick is a step for building the IaC: it can be a terraform directory, an ansible playbook, a kubernetes yaml file...
Most of those bricks depends of others. Actually, as for a construction, you can't put bricks of the second floor before building the first floor. And same, you can't add a user to a database before popping the database instance on the cloud provider.


### Common infra codes problem

In a lot of company, the infra code division in bricks and the right to apply these bricks were often strongly linked to technologies and the dependencies between bricks are not clear.
In a complexe infra code it leads to :
- mistakes from infra team when applying a change (because they don't understand dependencies) and they often guess how it should be applied.
- a big documentation not always up-to-date for each detail and specificity (as the dependency tree is not clear)
- shortcuts where you copy paste a generated value from an other brick (because conventions about informations exchanges between bricks don't exist)
- some infra code omissions because it can't be done (or it's to complicated) by terraform
- interrogations without answer when you discover an instance on your webUI cloud provider and you don't know where it has been coded.
- reapply regurlarly all the infra to avoid drift and CD complexity. It makes development less fluent. Moreover, sometimes, the CD failed due to a teammate modification that you have to correct before seeing your change applied on your brick (a bit annoying when you try to fix a production incident)

ExeIaC tries to address those problems.


### ExeIaC philosophy

The main philosophy of exeIaC is to dissociate the brick technology complexity from the way bricks interact with each other. It tries to respond to the *dependency inversion principle* of SOLID.

How to dissociate? More than a script that manage your infra code, exeIaC is a convention and each brick has to respect this convention. In fact, the brick contract is to be able to be applied, to display output and to display dependencies bricks... (exact convention presented further).

What brings this dessociation? Each brick can be applied in a different way with different technologies. However, we don't recommend that you use too many different way to apply your bricks (ExeIaC implement modules to factorize application code for brick of the same type), but it permits you some eventual shortcuts.

For example, if you need to open an ssh tunnel or switch to another vpn before running your terraform you can do it easily and you won't need to document that tricky specificity because it will be written clearly in the code.

In the same way, if you have to do something by hand, you can write it in your brick instead of executing a terraform apply. So, brick will get information from other bricks, displays some instructions to do by hand and then displays outputs to the dependencies. That type of shortcuts have to be considered ugly. However they're still better than write nothing, because it won't be invisibilized, it won't break any dependencies, it will warn anyone trying to apply the brick. Moreover, in most case it can be scripted and you won't need to display instruction to do by hand.

This dissociation will permit to reduce the documentation size and apply the rule of "no details or specificities written in the documentation".
Why this rule? Because documentation is hard to maintain up-to-date and exhaustive. And a short general up-to-date documentation will always be better than an attempt of exhaustivity where you don't know how to search something, where you're not sure that it exists and where you're not sure that it is still up-to-date or if it exist.
How this decoupling will allow to reduce documentation? Because you will document your brick inside your brick directory (where the code is written). The detail documentation will be easy to find as your tree structure should be intuitive. And when update the brick it is to update a README.md file if needed and to block the PR until it has been done. Moreover, when you debug, the documentation will be in the same directory!

As everything is coded (thanks for allowing easy shortcut) and the dependencies between bricks is clear the mistakes by applying something are reduce and the impact of a fail is well known.

You can also avoid to discover cloud provider resources that you can't link to a code by following a best practice: the path of your brick has to be tagged on the resources. 


### IaC vocabulary

- brick: a part of the infra that can be managed unitary. It can be the configuration of a monitoring user on a database or the creation a an entire cluster. A brick as a part of a wall is on top of other bricks that can be considered as dependencies and some brick can be on top of it that are the dependents. As we are in infra as code a brick is describe by code and the code can be execute. By convention we can execute some task on a brick : init, apply, output, destroy, validate, fmt, show_dependencies
- super brick: a super brick is a brick that contains other bricks it can be applied by apply sequentially all subricks it contains by following the dependency tree recursively.
- elementary brick: a brick that is not a super brick
- room: a super brick that is subjective cut of the infra that make sense. Generally a team have a specific ownership on a room and it can be assimilate to a domain (in DDD overview). But in a simple code overview lets say that a room is a git repository.
- other brick regroupment :
  - floor you can talk about floor the first floor is the room infra transverse to all domains. on the second floor you have the room client, backoffice, AI... it's the floor used by developers and feature teams. The third floor is transverse to all domains and manage users access on the infra.
the floor abstraction should simplify the dependencies abstraction.
  - column 
- furniture: a set of scripts to use in daily usage but that isn't use for building your infra. It can be a script to dump a databae. It is a good idea to put those sort of script on your infra as code so when you have to debug you have all your tools and documentation in the same place but obviously they are completely optionnal and they have nothing to do with dependency tree.
brick path: the full path of a brick on your computer
brick name: the name of the repo and the relative path from nome of repo to the brick
brick name sanitized: if you use brick name as tag on your cloud resources you may have to change slash 
 to underscore or dash. it's up to you to define a the sanitized regex

---

## HOWTO

### Install

As it is written in bash you just need to get code and precise where all part can be found. here an example :

- git clone repository `git clone git@github.com:arthur91f/exeiac.git`
- add an alias to your .bashrc `alias exeiac="$HOME/git/exeiac/exeiac.sh"`
- Create in your home a .exeiac file:
```exeiac_lib_path="$HOME/git/exeiac"
modules_path="$HOME/git/my-modules/modules"
room_paths_list="$HOME/git/my-infra-core
$HOME/git/my-team1-app
$HOME/git/my-team2-app
$HOME/git/my-infra-users-access"
function sanitize_function {
    sed -e "s|/[0-9]*-|--|g" -e 's|[^0-9a-zA-Z_-]|?|g' <<<"$1"
}
```

What are all those variables:

- exeiac_lib_path: exeiac script have been cut in multiple files containing bash functions. They are sourced at the start of the execution.
- room_paths_lits: list of your infra code repository (they need to contain directly access to superbrick)
- modules_path: it's a set of functions developed by you that you will use to apply your different infra code brick. See below for best practice and what function you have to implement ; but, basically you will have one module by technology you use (terraform, ansible) and you will specify how to do a "terraform plan,apply,output...", an "ansible plan,apply,output...", ... Only put modules in that directory as exeiac will try to source each file in that directory.
- sanitize_function: used to sanitized brick name if necessary.

### Create your modules

First how your module will be called:

- Your modules will overload default_module.sh functions so you don't have to rewrite everythings. Actually you also have access to all exeiac functions, but just don't overload them.
- Before executing your module action, exeiac will place you in the brick directory
- All options passed to exeiac will be passed to your module

Then remember few important things:

- Don't forget to overwrite is_brick_using_this_module function (return 0 is true, return 1 or an other int is false)
- Don't forget to support --non-interactive option or to code non-interactive functions
- show_dependents haven't to be implemented in the module (it's an exeiac function that use show_dependencies)

Here, a list of useful functions you can use:

- source "$modules_path/MODULE_NAME"
    to import all functions of an other module to use them
- copy_function SOURCE_NAME NEW_NAME
    so you can rename the apply function of your terraform module as terraform_apply and then overwrite your apply function
- import_modules_functions MODULE_NAME
    to do import and rename all MODULE function
- get_brick_sanitized_name
- get_arg -boolean=OPT1 "$@"  check if --OPT1 or -OPT1 is present in arguments
- get_arg -string=STR1 "$@"   check if STR1 option exist and display its value -STR1=VALUE -STR1 "VALUE"
- get_absolute_path
- merge_string_on_new_line

Then ... I think all have been said. Check our *Best Practice* chapter and *Short Introduction to Bash* if needed.

### Create your exeiac rooms

A room is a directory that contains bricks. exeiac will recognize files or directorys as bricks if their names begins by a number and a dash. ExeIaC will ignore bricks that aren't directly in the room's directory or directly inside a (super)brick directory.

The priority order of execution will correspond to the alphabetical order of path from your room.

**EXAMPLE:**
In *italic* directory that won't be executed
```
room-directory/
  1-init/
    1-monitoring.sh
    1-production.sh
    1-staging.sh
  2-envs/
    *README.md*
    1-production/
      1-network/
      2-ssh_server/
      2-cluster_k8s/
      3-prometheus/
      4-app/
        1-database/
          *how_to/*
            *1-dump_db*
            *2-restore_dump*
        2-k8s_deployment/
    1-staging/
      1-network/
      2-ssh_server/
      2-cluster_k8s/
      3-prometheus/
      *sav/*
        *4-app/*
          *1-database/*
            *how_to/*
              *1-dump_db*
              *2-restore_dump*
          *2-k8s_deployment/*
          *3-dns/*
          *4-access/*
    2-monitoring/
      1-network/
      2-ssh_server/
      3-prometheus_federated/
      4-grafana/
  *documentation/*
    *1-install_terraform*
    *2-on-call_duty_calendar*
```

Here the priority order for applying will be:
```
1-init/1-monitoring.sh
1-init/1-production.sh
1-init/1-staging.sh
2-envs/1-production/1-network/
2-envs/1-production/2-cluster_k8s/
2-envs/1-production/2-ssh_server/
2-envs/1-production/3-prometheus/
2-envs/1-production/4-app/1-database/
2-envs/1-production/4-app/2-k8s_deployment/
2-envs/1-staging/1-network/
...
```

### Now lets use exeiac command line

exeiac ACTION BRICK [OPTIONS]

- **BRICK**: directory or script that correspond to infra code. It can be the brickname (path beginning by the room (infra code git repository) name or an absolute or relative path.

- **ACTION**:
  - init: install or check that all tools are installed (terraform) and download dependencies as `terraform init`
  - apply: apply the brick codes
  - plan: a dry run of what it will be applied
  - output: display outputs that can be used by dependents bricks. Generated secrets or id.
  - destroy: revert the apply
  - help: display a the specific bricks help (usually it will always be the same as the basic help (all bricks should have the same behaviour except maybe for debugging)
  - validate: check that the code can be run (for CI)
  - fmt: a linter for the brick code
  - show_dependencies: get brick name whom output is needed or that have to be applied before
  - show_dependents: get all brick name that call the brick's output
  - show_dependencies_recursively
  - show_dependents_recursively
  - list_bricks: list all elementary brick of the super brick or all elementary bricks

- **OPTIONS**:
  - bricks-list BRICKS_LIST: we can also pass a bricks list instead of an only brick
  - bricks-specifier=SPECIFIER: we can use specifiers to not execute the action on the brick provided in argument but on it's dependencies or dependents according to the specifier provided.
  - non-interactive: the brick action will be executed without waiting answer from user
  - ... (use exeiac help to get more options)

- **SPECIFIER**: they can be combined with + as *"dependencies+dependents"*
  - all
  - dependencies
  - recursive_dependencies
  - dependents
  - recursive_dependents
  - dependents_dependencies
  - recursive_dependents_dependencies
  - selected(default)

---

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

---

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
4. Show dependencies should as much as possible present elementary_brick. It's more precise and it will drastically improve excution of exeiac with special options.

---

## DISCUSSION

### Why bash

For the first version, the aim is to execute bash instruction write in console so use bash should be the more simple way to do that to conserve coloration output of tools etc.
Moreover, even if bash isn't the easiest language to develop, it's not a such big deal for me, and it will be more simple to debug module for others as you just have to copy paste code on terminal.
Modules can call bash exeiac functions if needed. But better try to call them with exeiac. 


### Performance aim:

- exeIaC runtime to execute one bricks should be instant and very cheap on resources
- exeIaC runtime to plan its execute plan should be less than 3 seconds
- exeIaC should not open to much bash in the same time without parallelism
- enable a parallelism mode in future.

