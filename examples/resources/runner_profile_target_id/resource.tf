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

resource "waypoint_runner_profile" "target_id" {
  profile_name     = "example-newest-a"
  oci_url          = "hashicorp/waypoint-odr:latest"
  plugin_type      = "docker"
  default          = true
#  target_runner_id = "01G5GNJEYC7RVJNXFGMHD0HCDT"
  target_runner_labels = {
    app = "payments"
  }

  environment_variables = {
    VAULT_ADDR           = "https://localhost:8200"
    VAULT_CLIENT_TIMEOUT = "30s"
  }
}

output "env_vars" {
  value = waypoint_runner_profile.target_id.environment_variables
}