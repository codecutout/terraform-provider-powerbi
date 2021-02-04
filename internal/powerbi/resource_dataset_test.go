package powerbi

import (
	"fmt"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataset_basic(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
	var datasetID string
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
				resource "powerbi_dataset" "test" {
					workspace_id = powerbi_workspace.test.id
					default_mode = "push"
					name = "Acceptance Test Dataset"

					table {
						name = "entries"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "entryIndex"
							data_type = "int64"
						}
						column {
							name = "entryValue"
							data_type = "decimal"
						}
						column {
							name = "entryDate"
							data_type = "datetime"
						}
						column {
							name = "entrySuccessful"
							data_type = "bool"
						}

						measure {
							name = "sum of values"
							expression = "SUM([entryValue])"
						}
					}

					table {
						name = "entries-audit"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "modifiedBy"
							data_type = "string"
						}
					}

					relationship {
						name = "entries to entires-audit"
						from_table = "entries"
						from_column = "entryId"
						to_table = "entries-audit"
						to_column = "entryId"
						cross_filtering_behavior = "automatic"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerbi_dataset.test", "id"),
					testDatsetExistsWithName("powerbi_dataset.test", "Acceptance Test Dataset"),
					testPushDataSuccessful("powerbi_dataset.test", "entries", []map[string]interface{}{
						{
							"entryId":         "abc",
							"entryIndex":      3,
							"entryValue":      13.4,
							"entryDate":       "2000-01-01 12:34:56",
							"entrySuccessful": true,
						},
					}),
					set("powerbi_dataset.test", "id", &datasetID),
				),
			},

			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				resource "powerbi_dataset" "test" {
					workspace_id = powerbi_workspace.test.id
					default_mode = "push"
					name = "Acceptance Test Dataset"

					table {
						name = "entries"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "entryIndex"
							data_type = "int64"
						}
						column {
							name = "entryValue"
							data_type = "decimal"
						}
						column {
							name = "entryDate"
							data_type = "datetime"
						}

						# removed column
						# column {
						# 	name = "entrySuccessful"
						# 	data_type = "bool"
						# }

						measure {
							name = "sum of values"
							expression = "SUM([entryValue])"
						}
					}

					table {
						name = "entries-audit"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "modifiedBy"
							data_type = "string"
						}

						# New column
						column {
							name = "modifiedDate"
							data_type = "datetime"
						}
					}

					relationship {
						name = "entries to entires-audit"
						from_table = "entries"
						from_column = "entryId"
						to_table = "entries-audit"
						to_column = "entryId"
						cross_filtering_behavior = "automatic"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerbi_dataset.test", "id"),
					resource.TestCheckResourceAttrPtr("powerbi_dataset.test", "id", &datasetID),
					testDatsetExistsWithName("powerbi_dataset.test", "Acceptance Test Dataset"),
					testPushDataSuccessful("powerbi_dataset.test", "entries", []map[string]interface{}{
						{
							"entryId":    "abc",
							"entryIndex": 3,
							"entryValue": 13.4,
							"entryDate":  "2000-01-01 12:34:56",
						},
					}),
					testPushDataSuccessful("powerbi_dataset.test", "entries-audit", []map[string]interface{}{
						{
							"entryId":      "abc",
							"modifiedBy":   "joe",
							"modifiedDate": "2020-01-01T00:00:00",
						},
					}),
				),
			},

			// // final step checks importing the current state we reached in the step above
			// {
			// 	ResourceName:      "powerbi_workspace.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func TestAccDataset_forceNew(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
	var datasetID string
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
				resource "powerbi_dataset" "test" {
					workspace_id = powerbi_workspace.test.id
					default_mode = "push"
					name = "Acceptance Test Dataset"

					table {
						name = "entries"
						column {
							name = "entryId"
							data_type = "string"
						}
					}

					table {
						name = "entries-audit"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "modifiedBy"
							data_type = "string"
						}
					}

					relationship {
						name = "entries to entires-audit"
						from_table = "entries"
						from_column = "entryId"
						to_table = "entries-audit"
						to_column = "entryId"
						cross_filtering_behavior = "automatic"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerbi_dataset.test", "id"),
					testDatsetExistsWithName("powerbi_dataset.test", "Acceptance Test Dataset"),
					set("powerbi_dataset.test", "id", &datasetID),
				),
			},

			// second step adds new table
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				resource "powerbi_dataset" "test" {
					workspace_id = powerbi_workspace.test.id
					default_mode = "push"
					name = "Acceptance Test Dataset"

					table {
						name = "entries"
						column {
							name = "entryId"
							data_type = "string"
						}
					}

					table {
						name = "entries-audit"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "modifiedBy"
							data_type = "string"
						}
					}

					# new table
					table {
						name = "entries-description"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "description"
							data_type = "string"
						}
					}

					relationship {
						name = "entries to entires-audit"
						from_table = "entries"
						from_column = "entryId"
						to_table = "entries-audit"
						to_column = "entryId"
						cross_filtering_behavior = "automatic"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrNotEquals("powerbi_dataset.test", "id", &datasetID),
					set("powerbi_dataset.test", "id", &datasetID),
				),
			},

			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				resource "powerbi_dataset" "test" {
					workspace_id = powerbi_workspace.test.id
					default_mode = "push"
					name = "Acceptance Test Dataset"

					table {
						name = "entries"
						column {
							name = "entryId"
							data_type = "string"
						}
					}

					table {
						name = "entries-audit"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "modifiedBy"
							data_type = "string"
						}
					}

					
					table {
						name = "entries-description"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "description"
							data_type = "string"
						}
					}

					# updated relationship
					relationship {
						name = "entries to entires-description"
						from_table = "entries"
						from_column = "entryId"
						to_table = "entries-description"
						to_column = "entryId"
						cross_filtering_behavior = "automatic"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrNotEquals("powerbi_dataset.test", "id", &datasetID),
					set("powerbi_dataset.test", "id", &datasetID),
				),
			},

			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				resource "powerbi_dataset" "test" {
					workspace_id = powerbi_workspace.test.id
					default_mode = "push"
					name = "Acceptance Test Dataset"

					table {
						name = "entries"
						column {
							name = "entryId"
							data_type = "string"
						}
					}

					# Removed table
					# table {
					# 	name = "entries-audit"
					# 	column {
					# 		name = "entryId"
					# 		data_type = "string"
					# 	}
					# 	column {
					# 		name = "modifiedBy"
					# 		data_type = "string"
					# 	}
					# }
# 
					
					table {
						name = "entries-description"
						column {
							name = "entryId"
							data_type = "string"
						}
						column {
							name = "description"
							data_type = "string"
						}
					}

					# updated relationship
					relationship {
						name = "entries to entires-description"
						from_table = "entries"
						from_column = "entryId"
						to_table = "entries-description"
						to_column = "entryId"
						cross_filtering_behavior = "automatic"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrNotEquals("powerbi_dataset.test", "id", &datasetID),
					set("powerbi_dataset.test", "id", &datasetID),
				),
			},
		},
	})
}

func testDatsetExistsWithName(rn string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		workspaceID := rs.Primary.Attributes["workspace_id"]
		dataset, err := client.GetDatasetInGroup(workspaceID, rs.Primary.ID)
		if err != nil {
			return err
		}

		if dataset == nil {
			return fmt.Errorf("dataset with ID '%s' does not exist", rs.Primary.ID)
		}

		if expectedName != "" && dataset.Name != expectedName {
			return fmt.Errorf("dataset has name '%s' was expecting '%s'", dataset.Name, expectedName)
		}

		return nil
	}
}

func testPushDataSuccessful(rn string, tableName string, rows []map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		workspaceID := rs.Primary.Attributes["workspace_id"]
		datasetID := rs.Primary.ID

		err := client.PostRowsInGroup(workspaceID, datasetID, tableName, powerbiapi.PostRowsInGroupRequest{
			Rows: rows,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testCheckResourceAttrNotEquals(resourceName string, attrName string, test *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if *test == rs.Primary.Attributes[attrName] {
			return fmt.Errorf("resource property %s had value %s. Expecting it to not have this value", attrName, rs.Primary.Attributes[attrName])
		}
		return nil
	}
}
