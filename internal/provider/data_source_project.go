package provider

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Description: "A data source to read project configuration",
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Waypoint project",
			},
			"applications": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Applications associated with the Waypoint project",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"project_variables": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of variables in Key/value pairs associated with the Waypoint Project",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key of the variable",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "value of the variable",
						},
					},
				},
			},
			"data_source_git": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Configuration of Git repository where waypoint.hcl file is stored",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Url of git repository storing the waypoint.hcl file",
						},
						"git_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path in git repository when waypoint.hcl file is stored in a sub-directory",
						},
						"git_ref": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Git repository ref containing waypoint.hcl file",
						},
						"ignore_changes_outside_path": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Waypoint ignores changes outside path storing waypoint.hcl file",
						},
						"git_poll_interval_seconds": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval at which Waypoint should poll git repository for changes",
						},
						"file_change_signal": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates signal to be sent to any applications when their config files change.",
						},
					},
				},
			},
			"remote_runners_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Remote runners enabled for the project",
			},
			"app_status_poll_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Application status poll interval in seconds",
			},
			"git_auth_basic": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "Basic authentication details for Git",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Git username",
						},
						"password": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Git password",
						},
					},
				},
			},
			"git_auth_ssh": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "SSH authentication details for Git",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_user": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Git user associated with private key",
						},
						"passphrase": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Passphrase to use with private key",
						},
						"ssh_private_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Private key to authenticate to Git",
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	resourceProjectRead(context.TODO(), d, m)

	return diags
}

func flattenApplications(applications []*gen.Application) []interface{} {
	if applications != nil {
		apps := make([]interface{}, len(applications), len(applications))

		for a, application := range applications {
			app := make(map[string]interface{})

			app["name"] = application.Name
			apps[a] = app
		}
		return apps
	}
	return make([]interface{}, 0)
}

func flattenVariables(variables []*gen.Variable) []interface{} {
	if variables != nil {
		vars := make([]interface{}, len(variables), len(variables))

		for v, variable := range variables {
			vari := make(map[string]interface{})

			vari["name"] = variable.Name
			vari["value"] = variable.Value.(*gen.Variable_Str).Str
			vars[v] = vari
		}
		return vars
	}
	return make([]interface{}, 0)
}
