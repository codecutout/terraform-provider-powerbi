package powerbi

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"powerbi": testAccProvider,
	}
}

func TestProvider_validate(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{
		"POWERBI_TENANT_ID",
		"POWERBI_CLIENT_ID",
		"POWERBI_CLIENT_SECRET",
	}
	for _, requiredEnv := range requiredEnvs {
		if v := os.Getenv(requiredEnv); v == "" {
			t.Fatal(fmt.Sprintf("%s must be set for acceptance tests", requiredEnv))
		}
	}
}
