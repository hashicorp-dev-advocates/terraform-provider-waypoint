terraform {
  required_providers {
    waypoint = {
      source  = "local/hashicorp/waypoint"
      version = "0.1.0"
    }
  }
}

provider "waypoint" {
  waypoint_addr = "localhost:9701"
  token         = "HZCwuUtmrrpPuNwZEfzCLR6NunrkNukaMruVWSXKPAFRmf3ivYejV2PRnNPVCLKpRi8djai1QSNqJfFrZjPiBCY7SPsk7or1tgj64fcDavYPFDzFq3iykqEZgkcW9fJmzPgYNznhaZPjfffb9fm2AWavMsqbdwzBKL9e"
}

data "waypoint_project" "tf-test" {
  project_name = "tf-test"
}

output "tf_test_apps" {
  value = data.waypoint_project.tf-test.applications
}

output "tf_test_variables" {
  value = data.waypoint_project.tf-test.project_variables
}

output "tf_test_data_source_git" {
  value = data.waypoint_project.tf-test.data_source_git
}

output "app_status_poll" {
  value = data.waypoint_project.tf-test.app_status_poll_seconds
}

output "git_auth" {
  value = data.waypoint_project.tf-test.git_auth_basic
  sensitive = true
}

data "waypoint_project" "njack" {
  project_name = "njack"
}

output "njack_ssh" {
  value = data.waypoint_project.njack.git_auth_ssh
  sensitive = true
}

#output "url" {
#  value = data.waypoint_project.njack.data_source_git
#}