locals {
  network_ip_range    = var.network_ip_range
  private_env_domain  = var.private_domain_name
  internal_env_domain = var.internal_domain_name
  tag                 = "${var.env_name}/cluster_k8s"
  cluster_properties = {
    nodes = {
      master = {
        cpu    = 8
        memory = 8192
        disk   = 100
        roles  = "master"
        zone   = "eu-west1-a"
      }
      failover = {
        cpu    = 8
        memory = 8192
        disk   = 100
        roles  = "failover"
        zone   = "eu-west2-b"
      }
      reader1 = {
        cpu    = 8
        memory = 8192
        disk   = 100
        roles  = "reader"
        zone   = "eu-west1-b"
      }
      reader2 = {
        cpu    = 8
        memory = 8192
        disk   = 100
        roles  = "reader"
        zone   = "eu-west2-a"
      }
      reader_data = {
        cpu    = 4
        memory = 4096
        disk   = 100
        roles  = "reader"
        zone   = "eu-west1-a"
      }
    }
  }
  private_domain_name  = "cluster-main.${local.private_env_domain}"
  internal_domain_name = "cluster-main.${local.internal_env_domain}"
  instance_id          = "${var.project_id}/instance/${random_uuid.cluster.id}"
  private_ip           = cidrhost(local.network_ip_range, 3)
  public_ip            = "34.${join(".", random_integer.public_ip[*].id)}"
  template = {
    needs = {
      provider_project     = var.project_id
      network_id           = var.network_id
      network_ip_range     = local.network_ip_range
    }
    created = {
      tag                  = local.tag
      cluster_id           = local.instance_id
      private_ip           = local.private_ip
      load_balancer_ip     = local.public_ip
      private_domain_name  = local.private_domain_name
      internal_domain_name = local.internal_domain_name
      nodes                = random_uuid.nodes
      admin_creds          = {
        username = "${var.env_name}-admin"
        password = random_password.this
      }
      properties           = local.cluster_properties
    }
  }
}

resource "random_uuid" "cluster" {
}
resource "random_uuid" "nodes" {
  for_each = local.cluster_properties.nodes
  keepers = {
    cluster_id = random_uuid.cluster.id
  }
}

resource "random_password" "this" {
  length           = 16
  special          = true
}

resource "random_integer" "public_ip" {
  count = 3
  keepers = {
    instance_id = random_uuid.cluster.id
  }
  min = 0
  max = 255
}

resource "local_file" "this" {
  file_permission = "0644"
  content         = yamlencode(local.template)
  filename        = "${path.root}/CREATED_cluster_database.yml"
}
