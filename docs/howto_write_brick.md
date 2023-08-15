# Write your brick

## Generality

A brick is an element of your infra and it's IaC representation. It can be a 
terraform state, an ansible playbook, an helm charts, a combination all thoose 
things or even other things depending of your brick's module.

A brick is represented by a directory.

But a brick is not only the code describing the real infra element.
The elements describing a brick are :
- code: content of the directory
- outputs: the state of the brick, the value taken by things declared in codes
  as genereated password, given public ipaddress or just some value that can be 
  used by other bricks
- input: describe the context of the brick, it is composed by parts of outputs 
  of other bricks
- module events: 
- trigger: input change or other brick module's event that will trigger a lay or
  another action on the brick


## Directory convention: priority number

To re 


## Writing brick.yml file

