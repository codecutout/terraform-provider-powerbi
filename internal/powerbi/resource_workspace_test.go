package powerbi

import (
	"fmt"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccWorkspace_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", "Acceptance Test Workspace"),
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "name", "Acceptance Test Workspace"),
				),
			},
			// second step updates it with a new title
			{
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace - Updated"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", "Acceptance Test Workspace - Updated"),
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "name", "Acceptance Test Workspace - Updated"),
				),
			},
			// final step checks importing the current state we reached in the step above
			{
				ResourceName:      "powerbi_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkspace_skew(t *testing.T) {
	var workspaceID string
	config := `
	resource "powerbi_workspace" "test" {
		name = "Acceptance Test Workspace"
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_workspace.test", "id", &workspaceID),
				),
			},
			// second step skew new title
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateGroupAsAdmin(workspaceID, powerbiapi.UpdateGroupAsAdminRequest{
						Name: "Acceptance Test Workspace - Skewed",
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", "Acceptance Test Workspace"),
				),
			},
			// third step skew by deleting group
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.DeleteGroup(workspaceID)
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", "Acceptance Test Workspace"),
				),
			},
		},
	})
}

func testCheckWorkspaceExistsWithName(rn string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		workspace, err := client.GetGroup(rs.Primary.ID)
		if err != nil {
			return err
		}

		if workspace == nil {
			return fmt.Errorf("workspace with ID '%s' does not exist", rs.Primary.ID)
		}

		if expectedName != "" && workspace.Name != expectedName {
			return fmt.Errorf("workspace has name '%s' was expecting '%s'", workspace.Name, expectedName)
		}

		return nil
	}
}

func testAccCheckPowerbiWorkspaceDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	client := testAccProvider.Meta().(*powerbiapi.Client)

	// loop through the resources in state, verifying each widget
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerbi_workspace" {
			continue
		}

		// Retrieve our workspace by API lookup
		workspace, err := client.GetGroup(rs.Primary.ID)
		if err != nil {
			return err
		}
		if workspace != nil {
			return fmt.Errorf("workspace '%s' still exists", workspace.Name)
		}
	}

	return nil
}
