terraform {
  required_version = "1.3.6"
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.4.3"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.2.3"
    }
  }
}

variable "cloud_providers" {
    description = "map of cloud providers and their name, and organisation"
    type        = map(object({
      name = string
      organisation = string
      signup_url = string
    }))
}

variable "git_organisation" {
    description = "The cloud providers that correspond to git"
    type = object({
      name = string
      organisation = string
      signup_url = string
    })
}

variable "projects" {
    description = "List of projects"
    type = map(object({
      env = string
      cloud_provider = string
      project_id = string
    }))
}
