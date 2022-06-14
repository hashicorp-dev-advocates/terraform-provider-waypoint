package waypoint

import (
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccWaypointRunnerProfileTargetId(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRunnerProfileDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRunnerProfileTargetId(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "profile_name", rName),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "oci_url", "hashicorp/waypoint-odr:latest"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "plugin_type", "docker"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "default", "true"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "target_runner_id", "01G5GNJEYC7RVJNXFGMHD0HCDT"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "environment_variables.VAULT_ADDR", "https://localhost:8200"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "environment_variables.VAULT_CLIENT_TIMEOUT", "30s"),
				),
			},
		},
	})
}

func TestAccWaypointRunnerProfileTargetLabels(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRunnerProfileDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRunnerProfileTargetLabels(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "profile_name", rName),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "oci_url", "hashicorp/waypoint-odr:latest"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "plugin_type", "docker"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "default", "true"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "target_runner_labels.app", "payments"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "environment_variables.VAULT_ADDR", "https://localhost:8200"),
					resource.TestCheckResourceAttr(
						"waypoint_runner_profile.target_id", "environment_variables.VAULT_CLIENT_TIMEOUT", "30s"),
				),
			},
		},
	})
}

func testAccCheckRunnerProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runner_profile" {
			continue
		}

		// check that your destroy logic has executed if not return an error
	}

	return nil
}

func testAccRunnerProfileTargetId(name string) string {
	return fmt.Sprintf(`
resource "waypoint_runner_profile" "target_id" {
  profile_name     = "%s"
  oci_url          = "hashicorp/waypoint-odr:latest"
  plugin_type      = "docker"
  default          = true
  target_runner_id = "01G5GNJEYC7RVJNXFGMHD0HCDT"

  environment_variables = {
    VAULT_ADDR           = "https://localhost:8200"
    VAULT_CLIENT_TIMEOUT = "30s"
  }
}`, name)
}

func testAccRunnerProfileTargetLabels(name string) string {
	return fmt.Sprintf(`
resource "waypoint_runner_profile" "target_id" {
  profile_name     = "%s"
  oci_url          = "hashicorp/waypoint-odr:latest"
  plugin_type      = "docker"
  default          = true

  target_runner_labels = {
    app = "payments"
  }

  environment_variables = {
    VAULT_ADDR           = "https://localhost:8200"
    VAULT_CLIENT_TIMEOUT = "30s"
  }
}`, name)
}
