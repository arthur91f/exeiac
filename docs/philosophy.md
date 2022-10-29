# Philosophy

  - common infra code problems we try to solve with exeiac
  - define the brick concept
  - vocabulary sum up
  - how to solves problems and best practices discussion

## Common infra code problems we try to solve with exeIaC

- You found an ip adress or some cred that have pop from nowhere in your infra 
  code and don't know if it's still usefull. How it happens ? Because you have 
  copy paste a terraform output to your ansible groups var one year before.
- You deploy a new feature by applying a terraform state and forget to reapply
  an other directory else where. How it happens ? You haven't write it in the
  documentation or you forget to read the docs.
- You have a very big documentation full of exception, that is not 100% up to
  date.
- You reapply regularly all the infra to avoid drift.
- You take some shortcut and do few things by hand because it is complicated
  to write it on code for some edge case.
- As a new comers on the team, you found the code but you don't know how to
  deploy it (what creds, what inventory, what env vars to use, what tools
  version) 

## Define the brick concept

### What we will try to do with the brick concept

We try to find one uniq concept to descrive every infra elements. As in Object 
Oriented Programming where everything is an object, I want to say in infra 
everything is a brick.

In a DevOps approach, we want to have the more comprehensive scope of the 
infra : Infra embrace all infra elements and processes to deploy a product 
online.

So we will define infra as set of bricks. In next parts we will define new 
facets of brick concept step by step.

#### Starting point: concretely what is a brick ?

It could be : 
- a terraform state for popping an instance
- an ansible playbook for configuring an instance
- a script that pop and configure things on your infra
- an helm chart to pop your application
- …

#### Bricks order

But we can have higher order bricks too. Thus a brick can also be a terraform 
state and an ansible playbook together (for example to pop and configure an ssh
server). But an higher order brick can also be all you production environment 
containing lot of bricks.

#### Representation vs reality

As an infra code concept a brick is of course a part of this infra code. But it
is also the real infra element that this code describe : is the code is 
deployed, deploying, has the real infra element have drifted...
We can call that the status of the brick.

#### A more comprehensive code

Generally infra code is more a description of a static infra than a behaviour.
That's why indempotence is generally a properpty of infra code. But if we want 
to really describe all "codes" we shouldn't forget the command to deploy that 
code. Sometimes it is intuitive (cd my_dir ; terraform apply), sometimes it is
less :
- you have to link some variables files to the description
- you have to link some inventory to a playbook
- you have to be in specified network and start a VPN
- ...
Anyway we think that to be really comprehensive the behaviour to execute your 
code belongs to the brick and deserve to be clearly written somewhere.

#### What we can do with a brick

As we want want to see brick concept as an OOP interface, what are its method ?
- deploy: that we call lay in reference to the wall building vocabulary
  (moreother deploy doesn't represent well the fact that you make a small change 
  on code and want to re-deploy only this small change)
- remove: when you remove an instance or a part of one configuration.
- plan: when you just want to dry-run and check what will change.
- output: display some resource that your code have generated (cloud instance 
  id, password, ip...)
- ... (surely many other things as download code dependencies, lint code ...)

Note: in the following we will sometimes say "execute a brick". It will means
to call one of those methods.

#### Link between bricks

You can't configure an instance before this instance have been poped on cloud.
You can't configure an application that needs a database before the database 
exist because your application brick needs the database dns name and creds.

So a brick fall within a sequence and have dependencies.

Lets thinks to the different types of dependencies :
- simple data dependencies: you need some of the dependency brick output to 
  execute your current brick. If the dependency output change you need to 
  re-deploy your brick.
- weak data dependencies: you just need some data of the dependency brick to 
  execute your current brick. But your brick doesn't truely use the data.
  For example you need credentials to an hop ssh server to deploy the 
  configuration of your database. But your database doesn't really need the ssh
  server to be useful.
- state (or strong) dependency: you need a whole brick have been deployed to
  execute your brick. For example: the brick that add human user ssh key on your
  ssh server need that the server have been poped. As it is not linked to a 
  data, in doubt, you need to re-deploy the dependents bricks each time your 
  dependency have changed. If the ssh server have been destroyed and recreated,
  you need to redeploy human user ssh keys on it.
  In fact we should consider that it is a shortcut. Actually we can represent
  those type of dependency by a data (the instance id, the timestamp of the 
  last creation...). Or if it's too tricky to represent it by a data output
  the best practice should be to merge these 2 bricks to let to the brick the
  intelligence of this tricky dependency.

But keep in mind that all bricks doesn't depends of each other. A brick of your
production environment doesn't depends of a brick of your staging environment.

Also as you can't lay brick of your 3rd floor if you haven't build your second 
floor, the dependencies shouldn't be circular.

#### Go further: particular higher order brick

- room: represent a functionnal unit. Ok that's the case for all bricks but it 
  can be useful for you to cut your code in big room: in practice it can be a 
  git repository or represent a domain.
- a floor: is a layer of all your infra. Define your floor can let you have a 
  very intuitive approach of what bricks depends of which in big lines. Here an
  example:
  - cloud initialisation
  - infra core (network, firewall, bastion...)
  - application layer
  - monitoring
  - human user acces
  Here your application layer can depend of an infra-core brick but an 
  infra-core brick can't depends of the application layer. That doesn't mean 
  that an application can't create a subnetwork and add firewall rules.
- a tower: is a set of bricks that belongs to different layer that are
  independent of the rest of the infra 

