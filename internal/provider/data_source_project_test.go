package waypoint

import (
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceWaypointProject(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProject(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "project_name", regexp.MustCompile(rName)),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "remote_runners_enabled", regexp.MustCompile("true")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "data_source_git.0.git_url", regexp.MustCompile("https://github.com/hashicorp/waypoint-examples")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "data_source_git.0.git_path", regexp.MustCompile("docker/go")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "data_source_git.0.git_ref", regexp.MustCompile("HEAD")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "data_source_git.0.file_change_signal", regexp.MustCompile("some-signal")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "data_source_git.0.git_poll_interval_seconds", regexp.MustCompile("")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "app_status_poll_seconds", regexp.MustCompile("12")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "project_variables.#", regexp.MustCompile("3")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "git_auth_basic.0.username", regexp.MustCompile("test")),
					resource.TestMatchResourceAttr(
						"data.waypoint_project.test", "git_auth_basic.0.password", regexp.MustCompile("test")),
				),
			},
		},
	})
}

func testAccDataSourceProject(name string) string {
	return fmt.Sprintf(`
resource "waypoint_project" "test" {

  project_name           = "%s"
  remote_runners_enabled = true
  
	data_source_git {
    git_url  = "https://github.com/hashicorp/waypoint-examples"
    git_path = "docker/go"
    git_ref  = "HEAD"
    file_change_signal = "some-signal"
    git_poll_interval_seconds = 90
  }

  project_variables = {
    name = "rob"
    job  = "dev-advocate"
    conference = "HashiConf EU 2022"
  }


  app_status_poll_seconds = 12

  git_auth_basic {
  	username = "test" 
    password = "test" 
  }
}

data "waypoint_project" "test" {
  project_name = waypoint_project.test.project_name
  depends_on = [waypoint_project.test]
}
`, name)
}
