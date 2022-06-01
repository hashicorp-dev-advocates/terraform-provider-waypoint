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

// provider that can be used to obtain a waypoint client for acceptance tests
// this is configured in the testAccPreCheck
var waypointProvider *schema.Provider

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
						"waypoint_project.test", "data_source_git.0.git_url", "https://github.com/hashicorp/waypoint-examples"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_path", "docker/go"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_ref", "HEAD"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.file_change_signal", "some-signal"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_poll_interval_seconds", "90"),
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

	waypointProvider = Provider()

	err := waypointProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
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

func testAccProjectBasic(name string) string {
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

  git_auth_basic {
  	username = "test" # Required
    password = "test" # Required
  }
}`, name)
}
