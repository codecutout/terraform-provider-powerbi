package powerbi

import (
	"fmt"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceWorkspace_basic(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
	var workspaceName = fmt.Sprintf("Acceptance Test Data Source Workspace %s - Basic", workspaceSuffix)

	provider := Provider()
	provider.Configure(terraform.NewResourceConfigRaw(nil))
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
