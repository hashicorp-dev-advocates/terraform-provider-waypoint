package waypoint

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerName = "waypoint"

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
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "app_status_poll_seconds", "12"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "project_variables.name", "rob"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "project_variables.job", "dev-advocate"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "project_variables.conference", "HashiConf EU 2022"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "git_auth_basic.0.username", "test"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "git_auth_basic.0.password", "test"),
				),
			},
		},
	})
}

func TestAccWaypointProjectSsh(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckProjectDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSsh(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "project_name", rName),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "remote_runners_enabled", "true"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_url", "ssh://github.com/hashicorp/waypoint-examples"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_path", "docker/go"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_ref", "HEAD"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.file_change_signal", "some-signal"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "data_source_git.0.git_poll_interval_seconds", "90"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "git_auth_ssh.0.git_user", "devops-rob"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "git_auth_ssh.0.passphrase", "test-password"),
					resource.TestCheckResourceAttr(
						"waypoint_project.test", "git_auth_ssh.0.ssh_private_key", "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQCjcGqTkOq0CR3rTx0ZSQSIdTrDrFAYl29611xN8aVgMQIWtDB/\nlD0W5TpKPuU9iaiG/sSn/VYt6EzN7Sr332jj7cyl2WrrHI6ujRswNy4HojMuqtfa\nb5FFDpRmCuvl35fge18OvoQTJELhhJ1EvJ5KUeZiuJ3u3YyMnxxXzLuKbQIDAQAB\nAoGAPrNDz7TKtaLBvaIuMaMXgBopHyQd3jFKbT/tg2Fu5kYm3PrnmCoQfZYXFKCo\nZUFIS/G1FBVWWGpD/MQ9tbYZkKpwuH+t2rGndMnLXiTC296/s9uix7gsjnT4Naci\n5N6EN9pVUBwQmGrYUTHFc58ThtelSiPARX7LSU2ibtJSv8ECQQDWBRrrAYmbCUN7\nra0DFT6SppaDtvvuKtb+mUeKbg0B8U4y4wCIK5GH8EyQSwUWcXnNBO05rlUPbifs\nDLv/u82lAkEAw39sTJ0KmJJyaChqvqAJ8guulKlgucQJ0Et9ppZyet9iVwNKX/aW\n9UlwGBMQdafQ36nd1QMEA8AbAw4D+hw/KQJBANJbHDUGQtk2hrSmZNoV5HXB9Uiq\n7v4N71k5ER8XwgM5yVGs2tX8dMM3RhnBEtQXXs9LW1uJZSOQcv7JGXNnhN0CQBZe\nnzrJAWxh3XtznHtBfsHWelyCYRIAj4rpCHCmaGUM6IjCVKFUawOYKp5mmAyObkUZ\nf8ue87emJLEdynC1CLkCQHduNjP1hemAGWrd6v8BHhE3kKtcK6KHsPvJR5dOfzbd\nHAqVePERhISfN6cwZt5p8B3/JUwSR8el66DF7Jm57BM=\n-----END RSA PRIVATE KEY-----\n"),
				),
			},
		},
	})
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

  app_status_poll_seconds = 12

  git_auth_basic {
  	username = "test" # Required
    password = "test" # Required
  }
}`, name)
}

func testAccProjectSsh(name string) string {
	return fmt.Sprintf(`
resource "waypoint_project" "test" {

  project_name           = "%s"
  remote_runners_enabled = true
  
	data_source_git {
    git_url  = "ssh://github.com/hashicorp/waypoint-examples"
    git_path = "docker/go"
    git_ref  = "HEAD"
    file_change_signal = "some-signal"
    git_poll_interval_seconds = 90
  }

  git_auth_ssh {
    git_user        = "devops-rob" 
    passphrase      = "test-password" 
    ssh_private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCjcGqTkOq0CR3rTx0ZSQSIdTrDrFAYl29611xN8aVgMQIWtDB/
lD0W5TpKPuU9iaiG/sSn/VYt6EzN7Sr332jj7cyl2WrrHI6ujRswNy4HojMuqtfa
b5FFDpRmCuvl35fge18OvoQTJELhhJ1EvJ5KUeZiuJ3u3YyMnxxXzLuKbQIDAQAB
AoGAPrNDz7TKtaLBvaIuMaMXgBopHyQd3jFKbT/tg2Fu5kYm3PrnmCoQfZYXFKCo
ZUFIS/G1FBVWWGpD/MQ9tbYZkKpwuH+t2rGndMnLXiTC296/s9uix7gsjnT4Naci
5N6EN9pVUBwQmGrYUTHFc58ThtelSiPARX7LSU2ibtJSv8ECQQDWBRrrAYmbCUN7
ra0DFT6SppaDtvvuKtb+mUeKbg0B8U4y4wCIK5GH8EyQSwUWcXnNBO05rlUPbifs
DLv/u82lAkEAw39sTJ0KmJJyaChqvqAJ8guulKlgucQJ0Et9ppZyet9iVwNKX/aW
9UlwGBMQdafQ36nd1QMEA8AbAw4D+hw/KQJBANJbHDUGQtk2hrSmZNoV5HXB9Uiq
7v4N71k5ER8XwgM5yVGs2tX8dMM3RhnBEtQXXs9LW1uJZSOQcv7JGXNnhN0CQBZe
nzrJAWxh3XtznHtBfsHWelyCYRIAj4rpCHCmaGUM6IjCVKFUawOYKp5mmAyObkUZ
f8ue87emJLEdynC1CLkCQHduNjP1hemAGWrd6v8BHhE3kKtcK6KHsPvJR5dOfzbd
HAqVePERhISfN6cwZt5p8B3/JUwSR8el66DF7Jm57BM=
-----END RSA PRIVATE KEY-----
EOF
  }
}`, name)
}
