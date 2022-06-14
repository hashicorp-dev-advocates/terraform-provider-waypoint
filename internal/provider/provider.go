package waypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TerraformProviderProductUserAgent is included in the User-Agent header for
// any API requests made by the provider.
const TerraformProviderProductUserAgent = "terraform-provider-waypoint"

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAYPOINT_TOKEN", nil),
				Description: "Waypoint token to authenticate to Waypoint server",
			},
			"waypoint_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAYPOINT_ADDR", nil),
				Description: "Waypoint server address",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"waypoint_project":        dataSourceProject(),
			"waypoint_runner_profile": dataSourceRunnerProfile(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"waypoint_project":        resourceProject(),
			"waypoint_runner_profile": resourceRunnerProfile(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := Config{
			Token:        d.Get("token").(string),
			WaypointAddr: d.Get("waypoint_addr").(string),
		}
		return config.Client()
	}

	return provider
}
