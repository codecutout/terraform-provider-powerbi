package powerbi

import (
	"context"
	"fmt"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceWorkspace_basic(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
	var workspaceName = fmt.Sprintf("Acceptance Test Data Source Workspace %s - Basic", workspaceSuffix)

	provider := Provider()
	provider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	client := provider.Meta().(*powerbiapi.Client)
	response, _ := client.CreateGroup(powerbiapi.CreateGroupRequest{
		Name: workspaceName,
	})
	workspaceID := response.ID
	defer client.DeleteGroup(workspaceID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "powerbi_workspace" "test" {
					name = "%s"
				}
				`, workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerbi_workspace.test", "name", workspaceName),
					resource.TestCheckResourceAttr("data.powerbi_workspace.test", "id", workspaceID),
				),
			},
		},
	})
}
