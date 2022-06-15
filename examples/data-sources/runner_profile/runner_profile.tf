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
  token         = "rHhYzVXQBcK6NWBuN2DAAHsofoLBjPDJXKFwnZKxU9Xh13ShJs31muEV5YyJ97JhPyLsYQok1vYhpYqhWQisKuGf8DxuBRM9qbfXbap789YsBYGKFmkXvpdaLSM83moC139XkJ1rV5PR9Nwyxk64gHGk9CMroYvgG"
}

data "waypoint_runner_profile" "test" {
  id = "01G5K3Z29H87VRVYSJVBGQF7AM"
}

output "profile_name" {
  value = data.waypoint_runner_profile.test.profile_name
}

output "runner_default_bool" {
  value = data.waypoint_runner_profile.test.default
}

output "oci_url" {
  value = data.waypoint_runner_profile.test.oci_url
}

output "plugin_type" {
  value = data.waypoint_runner_profile.test.plugin_type
}

output "environment_variables" {
  value = data.waypoint_runner_profile.test.environment_variables
}

output "target_labels" {
  value = data.waypoint_runner_profile.test.target_runner_labels
}