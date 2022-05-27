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