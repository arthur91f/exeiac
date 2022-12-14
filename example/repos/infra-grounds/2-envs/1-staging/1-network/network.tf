locals {
  network_ip = "10.3.0.0/20"
  network_id = "${var.project_id}/network/${random_uuid.network.id}"
  nat_ip     = "34.${join(".", random_integer.nat_ip[*].id)}"
  network_template = {
    needs = {
      provider_project     = var.project_id
    }
    created = {
      network_ip = local.network_ip
      network_id = local.network_id
      nat_ip     = local.nat_ip
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
