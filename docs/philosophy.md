# Philosophy

- common infra code problems we try to solve with exeiac
- define the brick concept
- vocabulary sum up
- how to solve problems and best practices discussion

## Common infra code problems we try to solve with exeIaC

- You found an ip address or some cred that have pop from nowhere in your infra 
  code and don't know if it's still useful. How does it happen? Because you have
  copy-pasted a terraform output to your ansible groups var one year before.
- You deploy a new feature by applying a terraform state and forget to reapply
  an other directory else where. How does it happen? You haven't written it in 
  the documentation or you forget to read the docs.
- You have a very big documentation full of exceptions that is not 100% up to
  date.
- You reapply regularly all the infra to avoid drift.
- You take some shortcut and do few things by hand because it is complicated
  to write it on code for some edge case.
- As a new comers on the team, you found the code, but you don't know how to
  deploy it (what creds, what inventory, what env vars to use, what tools
  version) 

## Define the brick concept

### What we will try to do with the brick concept

We try to find one single concept to descrive every infra elements. As in Object
Oriented Programming where everything is an object, I want to say in infra 
everything is a brick.

In a DevOps approach, we want to have the more comprehensive scope of the infra:
Infra embrace all infra elements and processes to deploy a product online.

So, we will define infra as set of bricks. In next parts we will define new 
facets of brick concept step by step.

#### Starting point: concretely what is a brick ?

It could be : 
- a terraform state for popping an instance
- an ansible playbook for configuring an instance
- a script that pop and configure things on your infra
- a helm chart to pop your application
- …

But we can have higher order bricks too. Thus, a brick can also be a terraform 
state and an ansible playbook together (for example to pop and configure an ssh
server). But an higher order brick can also be all your production environment 
containing a lot of bricks.

#### Representation vs reality

As an infra code concept a brick is of course a part of this infra code. But it
is also the real infra element that this code describe: is the code is deployed,
deploying, has the real infra element have drifted...
We can call that the status of the brick.

#### A more comprehensive code

Generally infra code is more a description of a static infra than a behaviour.
That's why indempotence is generally a property of infra code. But if we want 
to really describe all "codes" we shouldn't forget the command to deploy that 
code. Sometimes it is intuitive (cd my_dir ; terraform apply), sometimes it is
less:
- you have to link some variables files to the description
- you have to link some inventory to a playbook
- you have to be in specified network and start a VPN
- ...
Anyway we think that being really comprehensive the behaviour to execute your 
code belongs to the brick and deserve to be clearly written somewhere.

#### What we can do with a brick

As we want to see brick concept as an OOP interface, what are its method?
- deploy: that we call "lay" in reference to the wall building vocabulary
  (more other deploy doesn't represent well the fact that you make a small
  change on code and want to re-deploy only this small change)
- remove: when you remove an instance or a part of one configuration.
- plan: when you just want to dry-run and check what will change.
- output: display some resources that your code has generated (cloud instance 
  id, password, ip...)
- ... (surely many other things as download code dependencies, lint code ...)

Note: in the following we will sometimes say "execute a brick". It will mean to
call one of those methods.

#### Link between bricks

You can't configure an instance before this instance have been popped on cloud.
You can't configure an application that needs a database before the database 
exist because your application brick needs the database dns name and creds.

So a brick fall within a sequence and have dependencies.

Lets thinks to the different types of dependencies:
- simple data dependencies: you need some of the dependency brick output to 
  execute your current brick. If the dependency output change you need to 
  re-deploy your brick.
- weak data dependencies: you just need some datas of the dependency brick to 
  execute your current brick. But your brick doesn't truly use the data.
  For example, you need credentials to an hop ssh server to deploy the 
  configuration of your database. But your database doesn't really need the ssh
  server to be useful.
- state (or strong) dependency: you need a whole brick have been deployed to
  execute your brick. For example: the brick that adds human user ssh key on
  your ssh server need that the server has been popped. As it is not linked to a
  data, in doubt, you need to re-deploy the dependents bricks each time your 
  dependency have changed. If the ssh server have been destroyed and recreated,
  you need to redeploy human user ssh keys on it.
  In fact, we should consider that it is a shortcut. Actually we can represent
  that type of dependency by a data (the instance id, the time stamp of the 
  last creation...). Or if it's too tricky to represent it by a data output
  the best practice should be to merge these 2 bricks to let to the brick the
  intelligence of this tricky dependency.

But keep in mind that all bricks don't depends of each other. A brick of your
production environment don't depends of a brick of your staging environment.

Also as you can't lay brick of your 3rd floor if you haven't built your second 
floor, the dependencies shouldn't be circular.

*ExeIaC vocabulary : name brick B relatively to brick A : *
- *direct previous brick: * a brick B is a direct previous of brick A if A needs 
  data of B outputs to be lay. So B needs to be layed before A.
- *linked previous brick: * it's the recursive concept of direct previous.
  A brick B is a linked previous of brick A if B is a direct previous brick of A 
  or a direct previous of a direct previous of A or ...
- *direct next brick: * a brick B is a direct next of A if A is a direct previous
  of B ie. if B need data from A output to be layed.
- *linked next brick: * it's the recursive concept of direct next. It represent
  all the brick that can be impacted by a change to brick A.
- * independant bricks: * two bricks are independant between each other if they 
  are not linked previous or linked next.
#### Go further: particular higher order brick

- room: represent a functional unit. OK, that's the case for all bricks, but it 
  can be useful for you to cut your code in big room: in practice it can be a 
  git repository or represent a domain.
- a floor: is a layer of all your infra. Define your floor can let you have a 
  very intuitive approach of what bricks depends on which in big lines. Here an
  example:
  - cloud initialization
  - infra core (network, firewall, bastion...)
  - application layer
  - monitoring
  - human user access
  Here, your application layer can depend on an infra-core brick but an 
  infra-core brick can't depend on the application layer. That doesn't mean 
  that an application can't create a sub-network and add firewall rules.
- a tower: is a set of bricks that belongs to different layers that are
  independent of the rest of the infra 

