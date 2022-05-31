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
				Type:     schema.TypeList,
				Computed: true,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			//"git_url": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"git_path": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"git_ref": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"ignore_changes_outside_path": &schema.Schema{
			//	Type:     schema.TypeBool,
			//	Computed: true,
			//},
			//"data_source_poll_enabled": &schema.Schema{
			//	Type:     schema.TypeBool,
			//	Computed: true,
			//},
			//"data_source_poll_interval": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"remote_runners_enabled": &schema.Schema{
			//	Type:     schema.TypeBool,
			//	Computed: true,
			//},
			//"file_change_signal": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"project_variables": &schema.Schema{
			//	Type:     schema.TypeMap,
			//	Computed: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//},
			"data_source_git": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"git_path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"git_ref": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ignore_changes_outside_path": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"git_poll_interval_seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_change_signal": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"remote_runners_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"git_auth_basic": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": &schema.Schema{
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
			"git_auth_ssh": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_user": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"passphrase": &schema.Schema{
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"ssh_private_key": &schema.Schema{
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	//client := m.(*WaypointClient).conn
	//
	//projectName := d.Get("project_name").(string)
	//project, err := client.GetProject(context.TODO(), projectName)
	//if err != nil {
	//	return diag.Errorf("Error retrieving the %s project", projectName)
	//}
	//
	//applications := flattenApplications(project.Applications)
	//variables := flattenVariables(project.Variables)
	//d.SetId(project.Name)
	//
	//d.Set("applications", applications)
	//d.Set("git_url", project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Url)
	//d.Set("git_path", project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Path)
	//d.Set("git_ref", project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Ref)
	//d.Set("ignore_changes_outside_path",
	//	project.DataSource.Source.(*gen.Job_DataSource_Git).Git.IgnoreChangesOutsidePath)
	//d.Set("variables", variables)
	//d.Set("remote_runners_enabled", project.RemoteEnabled)
	//d.Set("file_change_signal", project.FileChangeSignal)
	//
	//d.Set("data_source_poll_enabled", project.DataSourcePoll.Enabled)
	//if project.DataSourcePoll.Enabled == true {
	//	d.Set("data_source_poll_interval", project.DataSourcePoll.Interval)
	//}

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
