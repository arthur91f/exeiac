locals {
  staging_template = {
    needs = {
        ci_provider         = var.git_organisation.organisation
        staging_projects    = var.projects.staging.project_id
        monitoring_projects = var.projects.monitoring.project_id
    }
    created = {
        cloud = {
            user_name         = "ci-staging"
            user_key          = random_password.staging_user
            project_access    = ["all/rw"]
            monitoring_access = []
        }
        ci = {
            env = var.projects.staging.env
            secrets = [
                {
                    name  = "TF_VAR_username"
                    value = "ci-staging" # local.staging_template.created.users.name
                },            
                {
                    name  = "TF_VAR_userkey"
                    value = random_password.staging_user # local.staging_template.created.users.key
                }
            ]
        }
    }
  }
}

resource "random_password" "staging_user" {
  length           = 16
  special          = false
}

resource "local_file" "staging" {
  file_permission = "0644"
  content         = yamlencode(local.staging_template)
  filename        = "${path.root}/CREATED_staging_project&ci.yml"
}

output "staging_user" {
  sensitive = true
  value = {
    username = local.staging_template.created.cloud.user_name
    userkey = local.staging_template.created.cloud.user_key
  }
}
