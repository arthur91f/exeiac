locals {
  production_user_name = "ci-production"
  production_user_key = random_password.production_user
  production_template = {
    needs = {
        ci_provider         = var.git_organisation.organisation
        production_projects    = var.projects.production.project_id
        monitoring_projects = var.projects.monitoring.project_id
    }
    created = {
        cloud = {
            user_name         = local.production_user_name
            user_key          = local.production_user_key
            project_access    = ["all/rw"]
            monitoring_access = []
        }
        ci = {
            env = var.projects.production.env
            secrets = [
                {
                    name  = "TF_VAR_username"
                    value = local.production_user_name
                },            
                {
                    name  = "TF_VAR_userkey"
                    value = local.production_user_key
                }
            ]
        }
    }
  }
}

resource "random_password" "production_user" {
  length           = 16
  special          = false
}

resource "local_file" "production" {
  file_permission = "0644"
  content         = yamlencode(local.production_template)
  filename        = "${path.root}/CREATED_production_project&ci.yml"
}

output "production_user" {
  sensitive = true
  value = {
    username = local.production_template.created.cloud.user_name
    userkey = local.production_template.created.cloud.user_key
  }
}
