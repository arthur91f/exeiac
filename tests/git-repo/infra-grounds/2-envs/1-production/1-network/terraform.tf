terraform {
  required_version = "1.2.7"
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

variable "environment" {
  description = "Provider project_name, username and credentials"
  type = object({
    name         = string
    project_name = string
    username     = string
    credentials  = string
  })
  sensitive = true
}

variable "rooms_paths_list" {
  description = "Rooms paths list"
  type        = map(string)
}

variable "state_tag" {
  description = "Tag or label to put on cloud resources"
  type        = string
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

# data "terraform_remote_state" "" {
#   backend = "local"
#   config = {
#     path = ""
#   }
# }

