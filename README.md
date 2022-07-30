# exeIaC

Execute infra code according to a bricks dependencies metaphore

## VOCABULARY
### IaC vocabulary
brick: a part of the infra that can be managed unitary. It can be the configuration of a monitoring user on a database or the creation a an entire cluster. A brick as a part of a wall is on top of other bricks that can be considered as dependencies and some brick can be on top of it that are the dependents. As we are in infra as code a brick is describe by code and the code can be execute. By convention we can execute some task on a brick :
- init
- plan
- apply
- output
- destroy
- validate
- fmt
- show_dependencies
super brick: a super brick is a brick that contains other bricks it can be applied by apply sequentially all subricks it contains by following the dependency tree recursively.
elementary brick: a brick that is not a super brick
room: a super brick that is subjective cut of the infra that make sense. Generally a team have a specific ownership on a room and it can be assimilate to a domain (in DDD overview). But in a simple code overview lets say that a room is a git repository.
other brick regroupment :
  floor you can talk about floor the first floor is the room infra transverse to all domains. on the second floor you have the room client, backoffice, AI... it's the floor used by developers and feature teams. The third floor is transverse to all domains and manage users accesson the infra.
the floor abstraction should simplify the dependencies abstraction.
  environment
furniture: a set of scripts to use in daily usage but that isn't use for building your infra. It can be a script to dump a databae. It is a good idea to put those sort of script on your infra as code so when you have to debug you have all your tools and documentation in the same place but obviously they are completely optionnal and they have nothing to do with dependency tree.
brick path: the full path of a brick on your computer
brick name: the name of the repo and the relative path from nome of repo to the brick
brick name sanitized: if you use brick name as tag on your cloud resources you may have to change slash 
 to underscore or dash. it's up to you to define a the sanitized regex

### exeIaC tool code vocabulary
module: a module is a 
action:
dependencies is the need that an other brick is already applied or some previous brick output.

## CONVENTION & BEST PRACTICES

### Convention
All modules should implement the same action describe up. The simple way to do that is to heritate from the default module (source $modules_path/default.sh) so if you haven't implement fmt, you won't have bugs.
All actions should be idempotent.
All dependencies have to be visible. The default module permit to track all dependencies commented like that #EXEIAC:depends:BRICK_NAME or //EXEIAC:depends:BRICK_NAME
Then choose one unique convention beside for all your infra code:
- Apply action is always: plan, ask confirmation, apply. All interactive action should understand --exeiac-opts=non-interactive to skip plan and ask confirmation.
- No action should be interactive (the code will be simplier if you haven't to implement non-interactive option)
All brick are directory or file with name begins by a figure and a dash (2-database). It permit to identify brick and brick content prioritize the apply of different brick in the same directory and to add directory and file that are not re

### Philosophy

#### Agility (permit shortcut)
exeIaC philosophy is to bring agility and permit IaC vison not be driven by technology.
We think that an ugly shortcut is always better than not write anything in code even if it is written in documentation. Why ? few lines about a specific details lost in a general documentation won't warn you before applying mistakes.
The most important for us is to always know how to apply something and what to apply or check before apply. So a stupid script well placed that get output of a brick and just display "add databse monitoring user to grafana: grafana mypassword http://grafana.mycompany.co" is far enough for a shortcut and won't break any dependencies and thanks to it you will easily know that you need to add user credentials to grafana. it's in infra code. I'm not saying the best practice is to have many small scripts that display actions to do manually. But if more than 70% of your code use the same module 25% use other modules and you have less than 5% of shortcut more or less mature it's ok. Theese 5% are technical debt but it can't be forgotten (it's written in the code) and the day you debug you exactly see what is done. Moreover what we say is that the fact your teraform provider can't open an ssh tunnel shouldn't be an infra code architecture problem but only a local technical problem. If you are familiar with SOLID principles you should think the dependency inversion principle.

#### Single source of truth
Everything used in infra code come from infra code. Never get an auto-generated secret manually and copy-paste it in others brick. Just overload your output function to get that output and present it. So the dependency between the two bricks will be clear.

#### Apply and check what is needed (no more, no less)
If you want to be sure your infra code don't drift, you will reapply all. It's generally long and make your CI break your working ryhtms by waiting it finishes useless test. But the worst the CI can fail due to a mistake of your colleagues or known issue that isn't linked to your merge request.
So use --exeiac-opts=plan-dependencies-before-apply-recursively

#### Security
You can see easily secrets dependencies and what to give to who or rethink your infra architecture.

### Best practices 
Don't multiply modules that use different technology. Try to keep the applying infra code easy and intuitive. You should know how to apply a brick without reading its module ( Because it is always the sames ). We tend to say that an average infra collaborator should not masterize more than 3 modules and that masterize one complexe module is too much for most of developers.
However you can use heritage to get multiple modules that use a different terraform version. If they keep the same logic and it can be clearer than one unique complexe module that trying to manage different binary version.
Don't hesitate to display warning message or asking confirmation when you implement a tricky shortcut.
Tag your cloud resources with the brick_name or sanitized_brick_name so when you will see something in cloud it will be very easy to find it in code.
All outputs have to present the same format like json. The more you will define convention about outputs the more pass on change on dependents will be easy after modifying a brick.

## DISCUSSION
### Why bash
### Performance aim:
- exeIaC runtime to execute one bricks should be instant and very cheap on resources
- exeIaC runtime to plan its execute plan should be less than 3 seconds
- exeIaC should not open to much bash in the same time
- enable parallelism

infra-main (2tf) + (6tf + 1ansible + 2kube)*n
infra-apps (10tf 2tfm + 2kube + 2helm)*n + 2script*n
infra-monitoring (2tf 2helm 2kube)*n + (2tf 1helm 1 script)
infra-users (2tf 2tfm 1ansible)*n 1script

