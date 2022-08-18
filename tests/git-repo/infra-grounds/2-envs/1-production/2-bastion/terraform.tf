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

data "terraform_remote_state" "this_network" {
  backend = "local"
  config = {
    path = "${var.rooms_paths_list.infra-grounds}/2-envs/1-${var.environment.name}/1-network/terraform.tfstate"
  }
}

output "bastion" {
  sensitive = true
  value = {
    instance_id          = local.instance_id
    private_ip           = local.private_ip
    public_ip            = local.public_ip
    private_domain_name  = local.private_domain_name
    internal_domain_name = local.internal_domain_name
  }
}

