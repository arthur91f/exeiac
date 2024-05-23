locals {
  monitoring_template = {
    needs = {
        ci_provider         = var.git_organisation.organisation
        monitoring_projects    = var.projects.monitoring.project_id
        monitoring_projects = var.projects.monitoring.project_id
    }
    created = {
        cloud = {
            user_name         = "ci-monitoring"
            user_key          = random_password.monitoring_user
            project_access    = ["all/rw"]
            monitoring_access = []
        }
        ci = {
            env = var.projects.monitoring.env
            secrets = [
                {
                    name  = "TF_VAR_username"
                    value = "ci-monitoring" # local.monitoring_template.created.users.name
                },            
                {
                    name  = "TF_VAR_userkey"
                    value = random_password.monitoring_user # local.monitoring_template.created.users.key
                }
            ]
        }
    }
  }
}

resource "random_password" "monitoring_user" {
  length           = 16
  special          = false
}

resource "local_file" "monitoring" {
  file_permission = "0644"
  content         = yamlencode(local.monitoring_template)
  filename        = "${path.root}/CREATED_monitoring_project&ci.yml"
}

output "monitoring_user" {
  sensitive = true
  value = {
    username = local.monitoring_template.created.cloud.user_name
    userkey = local.monitoring_template.created.cloud.user_key
  }
}
