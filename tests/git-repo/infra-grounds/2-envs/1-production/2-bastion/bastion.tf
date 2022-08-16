locals {
  network_ip_range   = data.terraform_remote_state.this_network.outputs.network.ip_range
  private_env_domain = data.terraform_remote_state.this_network.outputs.domain_name.private
  internal_env_domain = data.terraform_remote_state.this_network.outputs.domain_name.internal
  
  instance_properties = {
    cpu    = 2
    memory = 4096
    disk   = {
      system = 30
    }
  }
  private_domain_name  = "bastion.${local.private_env_domain}"
  internal_domain_name = "bastion.${local.internal_env_domain}"
  instance_id = "${var.environment.project_name}/instance/${random_uuid.this.id}"
  private_ip = cidrhost(local.network_ip_range,2)
  public_ip  = "34.${join(".", random_integer.public_ip[*].id)}"
  template = {
    needs = {
      provider_project     = var.environment.project_name
      provider_username    = var.environment.username
      provider_credentials = var.environment.credentials
      network_ip_range     = local.network_ip_range
    }
    created = {
      tag         = var.state_tag
      instance_id = local.instance_id
      private_ip  = local.private_ip
      public_ip   = local.public_ip
      private_domain_name  = local.private_domain_name
      internal_domain_name = local.internal_domain_name
      properties           = local.instance_properties
    }
  }
}

resource "random_uuid" "this" {
}

resource "random_integer" "public_ip" {
  count = 3
  keepers = {
    instance_id = random_uuid.this.id
  }
  min = 0
  max = 255
}

resource "local_file" "this" {
  file_permission = "0444"
  content         = yamlencode(local.template)
  filename        = "${path.root}/CREATED_bastion.yml"
}

