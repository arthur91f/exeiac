locals {
  public_domain_name   = "myproduct.com"
  internal_domain_name = "${var.env_name}.myp.tech"
  private_domain_name  = "${var.env_name}.private"
  dns_template = {
    needs = {
      provider_project     = var.project_id
    }
    created = {
      public_domain_name   = local.public_domain_name
      internal_domain_name = local.internal_domain_name
      private_domain_name  = local.private_domain_name
    }
  }
}

resource "local_file" "dns" {
  file_permission = "0644"
  content         = yamlencode(local.dns_template)
  filename        = "${path.root}/CREATED_dns.yml"
}

