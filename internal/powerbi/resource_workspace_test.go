package powerbi

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccWorkspace_basic(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", fmt.Sprintf("Acceptance Test Workspace %s", workspaceSuffix)),
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "name", fmt.Sprintf("Acceptance Test Workspace %s", workspaceSuffix)),
				),
			},
			// second step updates it with a new title
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s - Updated"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", fmt.Sprintf("Acceptance Test Workspace %s - Updated", workspaceSuffix)),
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "name", fmt.Sprintf("Acceptance Test Workspace %s - Updated", workspaceSuffix)),
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

func TestAccWorkspace_capacity(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
	isPremiumCapacity := os.Getenv("POWERBI_IS_PREMIUM")
	premiumCapacityID := os.Getenv("POWERBI_CAPACITY_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			switch strings.ToLower(isPremiumCapacity) {
			case "":
				t.Fatal("POWERBI_IS_PREMIUM must be set for assigning capacity id to workspace acceptance tests")
			case "true":
				if premiumCapacityID == "" {
					t.Fatal("POWERBI_CAPACITY_ID must be set when POWERBI_IS_PREMIUM is set to \"true\" for capacity id acceptance tests")
				}
			case "false":
				t.Skip("Assign capacity id acceptance tests skipped")
			default:
				t.Fatal("POWERBI_IS_PREMIUM must be set to either \"true\" or \"false\"")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", fmt.Sprintf("Acceptance Test Workspace %s", workspaceSuffix)),
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "name", fmt.Sprintf("Acceptance Test Workspace %s", workspaceSuffix)),
				),
			},
			// second step assigns capacity id
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
					capacity_id = "%s"
				}
				`, workspaceSuffix, premiumCapacityID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "capacity_id", premiumCapacityID),
				),
			},
			// third step unassigns capacity id
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
					resource.TestCheckResourceAttr("powerbi_workspace.test", "capacity_id", ""),
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
	workspaceSuffix := acctest.RandString(6)
	config := fmt.Sprintf(`
	resource "powerbi_workspace" "test" {
		name = "Acceptance Test Workspace %s"
	}
	`, workspaceSuffix)

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
						Name: fmt.Sprintf("Acceptance Test Workspace %s - Skewed", workspaceSuffix),
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", fmt.Sprintf("Acceptance Test Workspace %s", workspaceSuffix)),
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
					testCheckWorkspaceExistsWithName("powerbi_workspace.test", fmt.Sprintf("Acceptance Test Workspace %s", workspaceSuffix)),
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
