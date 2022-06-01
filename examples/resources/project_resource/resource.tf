terraform {
  required_providers {
    waypoint = {
      source  = "local/hashicorp/waypoint"
      version = "0.1.0"
    }
  }
}

provider "waypoint" {
  waypoint_addr = "localhost:9701" # Address of Waypoint server does not require transport protocol like http
#  token         = "..."
  # token can be set with the $WAYPOINT_TOKEN environment variable or set directly
}

resource "waypoint_project" "example" {

  # Name of Waypoint Project - string
  project_name           = "tf-test" # Required
  # Should remote runner be enabled  - bool
  remote_runners_enabled = true # Defaults to `false`
  # Data source stanza for git config - map of below keys
  data_source_git {
    git_url  = "https://github.com/hashicorp/waypoint-examples" # Must be ssh url if using `git_auth_ssh`
    git_path = "docker/go" # Path in repo for waypoint.hcl file
    git_ref  = "HEAD" # This can be a branch name, a tag name, or a fully qualified Git ref such as refs/pull/1014
    file_change_signal = "some-signal"
    git_poll_interval_seconds = 15
  }

  app_status_poll_seconds = 12


  # Input variables for the Waypoint project - map with k/v pairs
  project_variables = {
    name = "rob"
    job  = "dev-advocate"
    conference = "HashiConf EU 2022"
  }



  # Git auth basic example stanza - map of below keys
  # Must either be git_auth_basic or git_auth_ssh.
  # Used to configure credentials for Waypoint to authenticate to Git
  #
    git_auth_basic {
      username = "test" # Required
      password = "test" # Required
    }
# Git auth ssh example stanza - map of below keys
# ssh_private_key is required and should be PKCS#1 type
#
#  git_auth_ssh {
#    git_user        = "devops-rob" # Required
#    passphrase      = "test-password" # Optional
#    ssh_private_key = <<EOF
#-----BEGIN RSA PRIVATE KEY-----
#MIICXAIBAAKBgQCjcGqTkOq0CR3rTx0ZSQSIdTrDrFAYl29611xN8aVgMQIWtDB/
#lD0W5TpKPuU9iaiG/sSn/VYt6EzN7Sr332jj7cyl2WrrHI6ujRswNy4HojMuqtfa
#b5FFDpRmCuvl35fge18OvoQTJELhhJ1EvJ5KUeZiuJ3u3YyMnxxXzLuKbQIDAQAB
#AoGAPrNDz7TKtaLBvaIuMaMXgBopHyQd3jFKbT/tg2Fu5kYm3PrnmCoQfZYXFKCo
#ZUFIS/G1FBVWWGpD/MQ9tbYZkKpwuH+t2rGndMnLXiTC296/s9uix7gsjnT4Naci
#5N6EN9pVUBwQmGrYUTHFc58ThtelSiPARX7LSU2ibtJSv8ECQQDWBRrrAYmbCUN7
#ra0DFT6SppaDtvvuKtb+mUeKbg0B8U4y4wCIK5GH8EyQSwUWcXnNBO05rlUPbifs
#DLv/u82lAkEAw39sTJ0KmJJyaChqvqAJ8guulKlgucQJ0Et9ppZyet9iVwNKX/aW
#9UlwGBMQdafQ36nd1QMEA8AbAw4D+hw/KQJBANJbHDUGQtk2hrSmZNoV5HXB9Uiq
#7v4N71k5ER8XwgM5yVGs2tX8dMM3RhnBEtQXXs9LW1uJZSOQcv7JGXNnhN0CQBZe
#nzrJAWxh3XtznHtBfsHWelyCYRIAj4rpCHCmaGUM6IjCVKFUawOYKp5mmAyObkUZ
#f8ue87emJLEdynC1CLkCQHduNjP1hemAGWrd6v8BHhE3kKtcK6KHsPvJR5dOfzbd
#HAqVePERhISfN6cwZt5p8B3/JUwSR8el66DF7Jm57BM=
#-----END RSA PRIVATE KEY-----
#EOF
#  }
}