package waypoint

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunnerProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRunnerProfileRead,
		Description: "A data source to read waypoint runner profiles",
		Schema: map[string]*schema.Schema{
			"profile_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the runner profile",
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Computed ID of runner profile.",
			},
			"oci_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "oci_url is the OCI image that will be used to boot the on demand runner.",
			},
			"plugin_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plugin type for runner i.e docker / kubernetes / aws-ecs.",
			},
			"plugin_config": {
				Type:        schema.TypeString, // Under the hood the type is []byte
				Computed:    true,
				Description: "plugin config is the configuration for the plugin that is created. It is usually HCL and is decoded like the other plugins, and is plugin specific.",
			},
			"plugin_config_format": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "config format specifies the format of plugin_config.",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if this runner profile is the default for any new projects",
			},
			"target_runner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the target runner for this profile.",
			},
			"target_runner_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of labels on target runners",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"environment_variables": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Any env vars that should be exposed to the on demand runner.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRunnerProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	resourceRunnerProfileRead(context.TODO(), d, m)

	return diags
}
