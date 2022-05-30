package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerName = "waypoint"

var providerFactories = map[string]func() (*schema.Provider, error){
	providerName: func() (*schema.Provider, error) { return Provider(), nil },
}

func TestAccWaypointProjectBasic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckProjectDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "project_name", rName),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "remote_runners_enabled", "true"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.data_source_poll_interval", "90s"),
					testDataSourceGit(),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("WAYPOINT_TOKEN") == "" {
		t.Fatal("Please set the environment variable WAYPOINT_TOKEN")
	}

	if os.Getenv("WAYPOINT_ADDR") == "" {
		t.Fatal("Please set the environment variable WAYPOINT_ADDR")
	}
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "waypoint_project" {
			continue
		}

		// check that your destroy logic has executed if not return an error
	}

	return nil
}

func testDataSourceGit() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "waypoint_project" {
				continue
			}

			projName := rs.Primary.Attributes["project_name"]

			// fetch the project from waypoint
			conn := Provider().Meta().(*WaypointClient).conn
			proj, err := conn.GetProject(context.TODO(), projName)

			if err != nil {
				return err
			}

			dpi := rs.Primary.Attributes["data_source_git.data_source_poll_interval"]
			if proj.DataSourcePoll.Interval != dpi {
				return fmt.Errorf("Poll Interval not set")
			}

			return nil
		}

		return nil
	}
}

func testAccProjectBasic(name string) string {
	return fmt.Sprintf(`
resource "waypoint_project" "test" {

  project_name           = "%s"
  remote_runners_enabled = true
  
	data_source_git {
    data_source_git_url  = "https://github.com/hashicorp/waypoint-examples"
    data_source_git_path = "docker/go"
    data_source_git_ref  = "HEAD"
    file_change_signal = "some-signal"
    data_source_poll_interval = "90s"
  }

  project_variables = {
    name = "rob"
    job  = "dev-advocate"
    conference = "HashiConf EU 2022"
  }

  git_auth_basic {
  	username = "test" # Required
    password = "test" # Required
  }
}`, name)
}
