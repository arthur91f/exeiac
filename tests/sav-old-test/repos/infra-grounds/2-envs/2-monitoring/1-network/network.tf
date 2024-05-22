locals {
  network_ip = "10.1.0.0/20"
  network_id = "${var.project_id}/network/${random_uuid.network.id}"
  nat_ip     = "34.${join(".", random_integer.nat_ip[*].id)}"
  peerings   = {
    production = "${local.network_id}-${var.production_network}",
    staging = "${local.network_id}-${var.staging_network}"
  }
  network_template = {
    needs = {
      provider_project     = var.project_id
    }
    created = {
      network_ip = local.network_ip
      network_id = local.network_id
      nat_ip     = local.nat_ip
      peerings   = local.peerings
    }
  }
}

resource "random_uuid" "network" {
  keepers = {
    network_ip = local.network_ip
  }
}

resource "random_integer" "nat_ip" {
  count = 3
  keepers = {
    network_ip = local.network_ip
  }
  min = 0
  max = 255
}

resource "local_file" "network" {
  file_permission = "0644"
  content         = yamlencode(local.network_template)
  filename        = "${path.root}/CREATED_network.yml"
}
