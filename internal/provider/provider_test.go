package waypoint

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	providerName: func() (*schema.Provider, error) { return Provider(), nil },
}

//func TestProvider(t *testing.T) {
//	if err := New("dev")().InternalValidate(); err != nil {
//		t.Fatalf("err: %s", err)
//	}
//}

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
