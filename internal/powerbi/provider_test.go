package powerbi

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"powerbi": testAccProvider,
	}
}

func TestProvider_validate(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_apiThrottling(t *testing.T) {

	provider := Provider()
	provider.Configure(terraform.NewResourceConfigRaw(nil))
	client := provider.Meta().(*powerbiapi.Client)

	c := make(chan int)

	workerGroup := sync.WaitGroup{}
	for workerIndex := 0; workerIndex < 3; workerIndex++ {
		workerGroup.Add(1)
		go func(workerIndex int) {
			defer workerGroup.Done()
			for v := range c {

				t1 := time.Now()
				_, err := client.GetGroups("", 1, 0)
				t2 := time.Now()

				t.Logf("worker %d: Request %d took %s", workerIndex, v, t2.Sub(t1))
				if err != nil {

					t.Error(err)
					return
				}
			}
		}(workerIndex)
	}

	for i := 0; i < 120; i++ {
		c <- i
	}
	close(c)

	workerGroup.Wait()
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
