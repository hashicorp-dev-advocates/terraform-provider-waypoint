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

resource "waypoint_project" "example" {

  project_name           = "example"
  remote_runners_enabled = false

  data_source_git {
    git_url                   = "https://github.com/hashicorp/waypoint-examples"
    git_path                  = "docker/go"
    git_ref                   = "HEAD"
    file_change_signal        = "some-signal"
    git_poll_interval_seconds = 15
  }

  app_status_poll_seconds = 12

  project_variables = {
    name       = "devopsrob"
    job        = "dev-advocate"
    conference = "HashiConf EU 2022"
  }
}