package provider

import (
	"context"
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
		ReadContext:   resourceProjectRead,
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
						"git_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"git_path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"git_ref": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ignore_changes_outside_path": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"git_poll_interval_seconds": &schema.Schema{
							Type:     schema.TypeInt,
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

	wp := m.(*WaypointClient).conn

	projectConf := client.DefaultProjectConfig()

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
			Url:                      dataSourceSlice["git_url"].(string),
			Path:                     dataSourceSlice["git_path"].(string),
			IgnoreChangesOutsidePath: dataSourceSlice["ignore_changes_outside_path"].(bool),
			Ref:                      dataSourceSlice["git_ref"].(string),
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
			Url:                      dataSourceSlice["git_url"].(string),
			Path:                     dataSourceSlice["git_path"].(string),
			IgnoreChangesOutsidePath: dataSourceSlice["ignore_changes_outside_path"].(bool),
			Ref:                      dataSourceSlice["git_ref"].(string),
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
	projectConf.Name = d.Get("project_name").(string)
	projectConf.RemoteRunnersEnabled = d.Get("remote_runners_enabled").(bool)

	if fileChangeSignal, ok := dataSourceSlice["file_change_signal"].(string); ok {
		projectConf.FileChangeSignal = fileChangeSignal
	}

	if dataSourcePollInterval, ok := dataSourceSlice["git_poll_interval_seconds"].(int); ok {
		projectConf.GitPollInterval = time.Duration(dataSourcePollInterval) * time.Second
	}

	_, err := wp.UpsertProject(context.TODO(), projectConf, gitConfig, variableList)

	if err != nil {
		return diag.FromErr(err)
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*WaypointClient).conn

	projectName := d.Get("project_name").(string)
	project, err := client.GetProject(context.TODO(), projectName)
	if err != nil {
		return diag.Errorf("Error retrieving the %s project", projectName)
	}

	d.SetId(project.Name)

	d.Set("remote_runners_enabled", project.RemoteEnabled)

	applications := flattenApplications(project.Applications)
	d.Set("applications", applications)

	variables := flattenVariables(project.Variables)
	d.Set("project_variables", variables)

	dataSourceGitSlice := map[string]interface{}{}
	dataSourceGitSlice["git_url"] = project.DataSource.GetGit().Url
	dataSourceGitSlice["git_path"] = project.DataSource.GetGit().Path
	dataSourceGitSlice["git_ref"] = project.DataSource.GetGit().Ref
	dataSourceGitSlice["file_change_signal"] = project.FileChangeSignal

	dpi, _ := time.ParseDuration(project.DataSourcePoll.Interval)
	dataSourceGitSlice["git_poll_interval_seconds"] = dpi / time.Second

	d.Set("data_source_git", []interface{}{dataSourceGitSlice})

	return nil
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
