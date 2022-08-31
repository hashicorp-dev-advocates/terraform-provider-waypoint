package waypoint

import (
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccWaypointAuthMethodOidc(t *testing.T) {
	amName := sdkacctest.RandomWithPrefix(providerName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckAuthMethodOidcDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthMethodOidc(amName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "name", amName),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "client_id", "060d0801-6fv7-4b03-ad59-b6397e2hc4ad"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "client_secret", "lFBsQ~Gb2E7p3Q3jr9jrWWw6aIQ5L3tYJ3wzLc1q"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "discovery_url", "https://login.microsoftonline.com/<insert-tenant-ID-here>/v2.0"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "auds.#", "2"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "auds.0", "testers"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "auds.1", "devs"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "allowed_redirect_urls.#", "2"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "allowed_redirect_urls.0", "https://localhost:9702/auth/oidc-callback"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "allowed_redirect_urls.1", "http://localhost:9701/oidc/callback"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "accessor_selector", "mycompany in list.groups"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "claim_mappings.%", "1"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "claim_mappings.groups", "groups"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "list_claim_mappings.%", "1"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "list_claim_mappings.groups", "groups"),
					resource.TestCheckResourceAttr(
						"waypoint_auth_method_oidc.test", "signing_algs.#", "1"),
				),
			},
		},
	},
	)
}

func testAccCheckAuthMethodOidcDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "auth_method_oidc" {
			continue
		}

		// check that your destroy logic has executed if not return an error
	}

	return nil
}

func testAccAuthMethodOidc(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "waypoint_auth_method_oidc" "test" {
  name          = var.name
  display_name	= var.name
  client_id     = "060d0801-6fv7-4b03-ad59-b6397e2hc4ad"
  client_secret = "lFBsQ~Gb2E7p3Q3jr9jrWWw6aIQ5L3tYJ3wzLc1q"
  discovery_url = "https://login.microsoftonline.com/<insert-tenant-ID-here>/v2.0"

  allowed_redirect_urls = [
    "https://localhost:9702/auth/oidc-callback",
    "http://localhost:9701/oidc/callback"
  ]
  
  auds = [
    "testers",
	"devs"
  ]

  claim_mappings = {
   groups = "groups"
  }


  list_claim_mappings = {
   "groups" = "groups"
  }

  accessor_selector = "mycompany in list.groups"

  signing_algs = [
    "rsa256"
  ]

  discovery_ca_pem = [
    "cert1.crt"
  ]

}
`, name)
}
