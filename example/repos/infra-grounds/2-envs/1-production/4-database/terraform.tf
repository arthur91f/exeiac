terraform {
  required_version = "1.3.6"
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.3.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.2.3"
    }
  }
}

variable "project_id" {
  description = "provider project name"
  type = string
}

variable "env_name" {
  description = "environment name"
  type = string
}

variable "network_id" {
  description = "the bastion network id"
  type = string
}

variable "network_ip_range" {
  description = "the private ip range of the network"
  type = string
}

variable "private_domain_name" {
  description = "the bastion's private domain name"
  type = string
}

variable "internal_domain_name" {
  description = "the bastion's internal domain name"
  type = string
}

output "cluster" {
  sensitive = true
  value = {
    instance_id          = local.instance_id
    private_ip           = local.private_ip
    public_ip            = local.public_ip
    private_domain_name  = local.private_domain_name
    internal_domain_name = local.internal_domain_name
  }
}

