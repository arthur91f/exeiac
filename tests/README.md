# EXAMPLE

This directory is here to show an example of infra code that respect the exeiac
convention. It can inspire you to do something same

The second aim is to run test on it.

## Structure

If you have read the README.md and docs/ directory it should seems natural to 
you.

- a directory that contains all your git repository of your infra
- a directory that contains all your modules
- an exeiac conf file

├── exeiac.yml
├── README.md
└── repos
    ├── app-backend
    ├── app-frontend
    ├── infra-grounds
    │   ├── 1-init
    │   │   └── brick.yml
    │   └── 2-envs
    │       ├── 1-production
    │       │   ├── 1-network
    │       │   │   └── brick.yml
    │       │   └── 2-bastion
    │       │       └── brick.yml
    │       ├── 1-staging
    │       │   ├── 1-network
    │       │   │   └── brick.yml
    │       │   └── 2-bastion
    │       │       └── brick.yml
    │       └── 2-monitoring
    │           ├── 1-network
    │           │   └── brick.yml
    │           └── 2-bastion
    │               └── brick.yml
    ├── modules
    │   ├── module_test.sh
    │   └── terraform.sh
    └── users

