resource "waypoint_project" "example" {

  project_name           = "example"
  remote_runners_enabled = true

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