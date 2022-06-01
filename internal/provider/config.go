package waypoint

import (
	"fmt"
	//gofastly "github.com/fastly/go-fastly/v6/fastly"
	waypoint "github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

type Config struct {
	Token        string
	WaypointAddr string
}

type WaypointClient struct {
	conn waypoint.Waypoint
}

func (c *Config) Client() (*WaypointClient, diag.Diagnostics) {
	var client WaypointClient

	if c.Token == "" {
		return nil, diag.FromErr(fmt.Errorf("[Err] No Waypoint token set"))
	}

	waypointClientConfig := waypoint.DefaultConfig()
	waypointClientConfig.Address = c.WaypointAddr
	waypointClientConfig.Token = c.Token

	waypointClient, err := waypoint.New(waypointClientConfig)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.conn = waypointClient
	return &client, nil
}
