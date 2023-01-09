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
  sensitive = true
}

variable "env_name" {
  description = "provider project name"
  type = string
  sensitive = true
}

output "domain_name" {
  sensitive = true
  value = {
    public   = local.public_domain_name
    internal = local.internal_domain_name
    private  = local.private_domain_name
  }
}

output "network" {
  sensitive = true
  value = {
    ip_range   = local.network_ip
    network_id = local.network_id
    nat_ip     = local.nat_ip
  }
}
