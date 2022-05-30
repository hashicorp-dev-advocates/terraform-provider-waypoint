package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Project resource in the Waypoint Terraform provider.",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceScaffoldingRead,
		UpdateContext: resourceProjectCreate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Waypoint project",
			},
			"project_variables": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_source_git": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source_git_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"data_source_git_path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"data_source_git_ref": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"data_source_git_ignore_changes_outside_path": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"data_source_poll_interval": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"file_change_signal": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"remote_runners_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"git_auth_basic": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"password": &schema.Schema{
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"git_auth_ssh": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"git_auth_basic"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_user": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"passphrase": &schema.Schema{
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"ssh_private_key": &schema.Schema{
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	wp := m.(*WaypointClient).conn

	// Git configuration for Waypoint project
	var gitConfig *client.Git

	authBasicList := d.Get("git_auth_basic").([]interface{})
	authSshList := d.Get("git_auth_ssh").([]interface{})

	dataSourceList := d.Get("data_source_git").([]interface{})
	dataSourceSlice := dataSourceList[0].(map[string]interface{})

	if len(authBasicList) > 0 {
		var auth *client.GitAuthBasic

		authBasicSlice := authBasicList[0].(map[string]interface{})
		username := authBasicSlice["username"]
		password := authBasicSlice["password"]

		auth = &client.GitAuthBasic{
			Username: username.(string),
			Password: password.(string),
		}

		gitConfig = &client.Git{
			Url:                      dataSourceSlice["data_source_git_url"].(string),
			Path:                     dataSourceSlice["data_source_git_path"].(string),
			IgnoreChangesOutsidePath: dataSourceSlice["data_source_git_ignore_changes_outside_path"].(bool),
			Ref:                      dataSourceSlice["data_source_git_ref"].(string),
			Auth:                     auth,
		}
	} else if len(authSshList) > 0 {
		var auth *client.GitAuthSsh
		authSshSlice := authSshList[0].(map[string]interface{})
		var passphrase interface{}
		gitUser := authSshSlice["git_user"]
		sshPrivateKey := authSshSlice["ssh_private_key"]
		if authSshSlice["passphrase"] != nil {
			passphrase = authSshSlice["passphrase"]
		}

		auth = &client.GitAuthSsh{
			User:          gitUser.(string),
			PrivateKeyPem: []byte(sshPrivateKey.(string)),
			Password:      passphrase.(string),
		}

		gitConfig = &client.Git{
			Url:                      dataSourceSlice["data_source_git_url"].(string),
			Path:                     dataSourceSlice["data_source_git_path"].(string),
			IgnoreChangesOutsidePath: dataSourceSlice["data_source_git_ignore_changes_outside_path"].(bool),
			Ref:                      dataSourceSlice["data_source_git_ref"].(string),
			Auth:                     auth,
		}

	}

	// Project variables configuration
	var variableList []*gen.Variable
	varsList := d.Get("project_variables").(map[string]interface{})

	for key, value := range varsList {

		projectVariable := client.SetVariable()
		projectVariable.Name = key
		projectVariable.Value = &gen.Variable_Str{Str: value.(string)}
		variableList = append(variableList, &projectVariable)
	}

	// Project config for request
	projectName := d.Get("project_name").(string)
	d.SetId(projectName)
	projectConf := client.DefaultProjectConfig()
	projectConf.Name = d.Get("project_name").(string)
	projectConf.RemoteRunnersEnabled = d.Get("remote_runners_enabled").(bool)

	if fileChangeSignal, ok := dataSourceSlice["file_change_signal"].(string); ok {
		projectConf.FileChangeSignal = fileChangeSignal
	}

	if dataSourcePollInterval, ok := dataSourceSlice["data_source_poll_interval"].(string); !ok {
		gpi, err := time.ParseDuration(dataSourcePollInterval)
		if err != nil {
			return diag.FromErr(fmt.Errorf("please specify data_source_poll_interval as Go duration string. E.g 30s, 5m: %s", err))
		}

		projectConf.GitPollInterval = gpi
	}

	_, err := wp.UpsertProject(context.TODO(), projectConf, gitConfig, variableList)

	if err != nil {
		return diag.FromErr(err)
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	return diags
}

func resourceScaffoldingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*WaypointClient).conn

	projectName := d.Get("project_name").(string)
	project, err := client.GetProject(context.TODO(), projectName)
	if err != nil {
		return diag.Errorf("Error retrieving the %s project", projectName)
	}

	applications := flattenApplications(project.Applications)
	variables := flattenVariables(project.Variables)
	d.SetId(project.Name)

	d.Set("applications", applications)
	d.Set("data_source_git_url", project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Url)
	d.Set("data_source_git_path", project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Path)
	d.Set("data_source_git_ref", project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Ref)
	d.Set("data_source_git_ignore_changes_outside_path",
		project.DataSource.Source.(*gen.Job_DataSource_Git).Git.IgnoreChangesOutsidePath)
	d.Set("variables", variables)
	d.Set("remote_runners_enabled", project.RemoteEnabled)
	d.Set("file_change_signal", project.FileChangeSignal)

	d.Set("data_source_poll_enabled", project.DataSourcePoll.Enabled)
	if project.DataSourcePoll.Enabled == true {
		d.Set("data_source_poll_interval", project.DataSourcePoll.Interval)
	}

	dataSourceSlice := map[string]interface{}{}
	dataSourceSlice["data_source_poll_interval"] = project.DataSourcePoll.Interval
	d.Set("data_source_git", dataSourceSlice)

	return diags
}

//func resourceScaffoldingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	// use the meta value to retrieve your client from the provider configure method
//	// client := meta.(*apiClient)
//
//	return diag.Errorf("not implemented")
//}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	var diags diag.Diagnostics
	wp := m.(*WaypointClient).conn

	//auth := gen.GitAuthBasic{
	//	Username: "",
	//	Password: "",
	//}

	gc := client.Git{
		Url:                      "RESOURCE DELETED",
		Path:                     "RESOURCE DELETED",
		IgnoreChangesOutsidePath: true,
		Ref:                      "RESOURCE DELETED",
		Auth:                     nil,
	}

	projectName := d.Get("project_name").(string)
	d.SetId(projectName)
	projectConf := client.DefaultProjectConfig()
	projectConf.Name = d.Get("project_name").(string)

	_, err := wp.UpsertProject(context.TODO(), projectConf, &gc, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	return diags

}

func gitAuthBasicMap() {

}
