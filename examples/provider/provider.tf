terraform {
  required_providers {
    waypoint = {
      source  = "hashicorp-dev-advocates/waypoint"
      version = "0.2.1"
    }
  }
}

provider "waypoint" {
  waypoint_addr = "localhost:9701"
  token         = "..."
}