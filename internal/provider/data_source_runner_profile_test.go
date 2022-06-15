package waypoint

import (
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceRunnerProfileId(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRunnerProfileId(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "profile_name", regexp.MustCompile(rName)),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "oci_url", regexp.MustCompile("hashicorp/waypoint-odr:latest")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "plugin_type", regexp.MustCompile("docker")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "default", regexp.MustCompile("true")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "target_runner_id", regexp.MustCompile("01G5GNJEYC7RVJNXFGMHD0HCDT")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "environment_variables.VAULT_ADDR", regexp.MustCompile("https://localhost:8200")),
				),
			},
		},
	})
}

func testAccDataSourceRunnerProfileId(name string) string {
	return fmt.Sprintf(`
resource "waypoint_runner_profile" "target_id" {
  profile_name     = "%s"
  oci_url          = "hashicorp/waypoint-odr:latest"
  plugin_type      = "docker"
  default          = true
  target_runner_id = "01G5GNJEYC7RVJNXFGMHD0HCDT"

  environment_variables = {
    VAULT_ADDR           = "https://localhost:8200"
  }
}

data "waypoint_runner_profile" "target_id" {
  id = waypoint_runner_profile.target_id.id
  depends_on = [
	waypoint_runner_profile.target_id
  ]
}

`, name)
}

func TestAccDataSourceRunnerProfileLabels(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRunnerProfileLabels(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "profile_name", regexp.MustCompile(rName)),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "oci_url", regexp.MustCompile("hashicorp/waypoint-odr:latest")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "plugin_type", regexp.MustCompile("docker")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "default", regexp.MustCompile("true")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "target_runner_labels.app", regexp.MustCompile("payments")),
					resource.TestMatchResourceAttr(
						"data.waypoint_runner_profile.target_id", "environment_variables.VAULT_ADDR", regexp.MustCompile("https://localhost:8200")),
				),
			},
		},
	})
}

func testAccDataSourceRunnerProfileLabels(name string) string {
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
  }
}

data "waypoint_runner_profile" "target_id" {
  id = waypoint_runner_profile.target_id.id
  depends_on = [
	waypoint_runner_profile.target_id
  ]
}

`, name)
}
