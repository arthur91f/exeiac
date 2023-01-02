locals {
  user_name = "registry-user"
  user_key  = random_password.user
  url       = "${var.projects.monitoring.project_id}/container-registry"
  registry_template = {
    needs = {
        monitoring_projects = var.projects.monitoring.project_id
    }
    created = {
        registry = {
            url               = local.url
        }
        user = {
          username = local.user_name
          userkey  = local.user_key
        }
    }
  }
}

resource "random_password" "user" {
  length           = 16
  special          = false
}

resource "local_file" "registry" {
  file_permission = "0644"
  content         = yamlencode(local.registry_template)
  filename        = "${path.root}/CREATED_registry.yml"
}

output "user" {
  sensitive = true
  value = {
    username = local.registry_template.created.user.username
    userkey = local.registry_template.created.user.userkey
  }
}

output "registry" {
  value = {
    url = local.registry_template.created.registry.url
  }
}
