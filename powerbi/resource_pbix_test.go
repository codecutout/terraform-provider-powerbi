package powerbi

import (
	"fmt"
	"github.com/alex-davies/terraform-provider-powerbi/powerbi/internal/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccPBIX_basic(t *testing.T) {
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

				resource "powerbi_pbix" "test" {
					workspace = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					content_base64 = filebase64("./resource_pbix_test.pbix")
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckReportExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "id"),
					resource.TestCheckResourceAttr("powerbi_pbix.test", "name", "Acceptance Test PBIX"),
				),
			},
		},
	})
}

func testCheckDatasetExistsInWorkspace(workspaceResourceName string, expectedDatasetName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[workspaceResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", workspaceResourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*api.Client)
		datasets, err := client.GetDatasetsInGroup(api.GetDatasetsInGroupRequest{
			GroupID: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		var datasteNames []string
		for _, dataset := range datasets.Value {
			if dataset.Name == expectedDatasetName {
				return nil
			}
			datasteNames = append(datasteNames, dataset.Name)
		}
		return fmt.Errorf("workspace has datasets %v was expecting list to contain '%s'", datasteNames, expectedDatasetName)
	}
}

func testCheckReportExistsInWorkspace(workspaceResourceName string, expectedReportName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[workspaceResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", workspaceResourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*api.Client)
		reports, err := client.GetReportsInGroup(api.GetReportsInGroupRequest{
			GroupID: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		var reportNames []string
		for _, report := range reports.Value {
			if report.Name == expectedReportName {
				return nil
			}
			reportNames = append(reportNames, report.Name)
		}
		return fmt.Errorf("workspace has reports %v was expecting list to contain '%s'", reportNames, expectedReportName)
	}
}
