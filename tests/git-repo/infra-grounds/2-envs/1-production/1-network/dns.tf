locals {
  public_domain_name   = "myproduct.com"
  internal_domain_name = "${var.environment.name}.myp.tech"
  private_domain_name  = "${var.environment.name}.private"
  dns_template = {
    needs = {
      provider_project     = var.environment.project_name
      provider_username    = var.environment.username
      provider_credentials = var.environment.credentials
    }
    created = {
      state_tag            = var.state_tag
      public_domain_name   = local.public_domain_name
      internal_domain_name = local.internal_domain_name
      private_domain_name  = local.private_domain_name
    }
  }
}

resource "local_file" "dns" {
  file_permission = "0444"
  content         = yamlencode(local.dns_template)
  filename        = "${path.root}/CREATED_dns.yml"
}

