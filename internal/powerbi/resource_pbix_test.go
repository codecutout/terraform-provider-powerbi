package powerbi

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/codecutout/terraform-provider-powerbi/internal/pbixrewriter"
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPBIX_basic(t *testing.T) {
	var updatedTime time.Time
	pbixLocation := TempFileName("", ".pbix")
	pbixLocationTfFriendly := strings.ReplaceAll(pbixLocation, "\\", "\\\\")
	workspaceSuffix := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				PreConfig: func() {
					Copy("./resource_pbix_test_sample1.pbix", pbixLocation)
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
				}
				`, workspaceSuffix, pbixLocationTfFriendly, pbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					setUpdatedTime("powerbi_pbix.test", &updatedTime),
					testCheckDatasetExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckReportExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "id"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "dataset_id"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "report_id"),
					resource.TestCheckResourceAttr("powerbi_pbix.test", "name", "Acceptance Test PBIX"),
				),
			},
			// update with different pbix same path
			{
				PreConfig: func() {
					Copy("./resource_pbix_test_sample2.pbix", pbixLocation)
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
				}
				`, workspaceSuffix, pbixLocationTfFriendly, pbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					testCheckUpdatedAfter("powerbi_pbix.test", &updatedTime), //update has occured since creation
					testCheckDatasetExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckReportExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "id"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "dataset_id"),
					resource.TestCheckResourceAttrSet("powerbi_pbix.test", "report_id"),
					resource.TestCheckResourceAttr("powerbi_pbix.test", "name", "Acceptance Test PBIX"),
				),
			},
			// deletes the resource
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckReportDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckResourceRemoved("powerbi_pbix.test"),
				),
			},
		},
	})
}

func TestAccPBIX_external_dataset_report(t *testing.T) {
	datasetPbixLocation := TempFileName("", ".pbix")
	datasetPbixLocationTfFriendly := strings.ReplaceAll(datasetPbixLocation, "\\", "\\\\")

	reportPbixLocation := TempFileName("", ".pbix")
	reportPbixLocationTfFriendly := strings.ReplaceAll(reportPbixLocation, "\\", "\\\\")

	workspaceSuffix := acctest.RandString(6)
	var datasetID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step create a dataset
			{
				PreConfig: func() {
					Copy("./resource_pbix_dataset_only.pbix", datasetPbixLocation)
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}
				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetExistsInWorkspace("powerbi_workspace.test", "Acceptance Test dataset PBIX"),
					testCheckReportDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test dataset PBIX"),
					set("powerbi_pbix.dataset_only", "dataset_id", &datasetID),
				),
			},
			// second step creates a report that links to the dataset
			{
				PreConfig: func() {
					// Will do rewrite our pbix so it points to the dataset
					pbixrewriter.RewritePbixFiles("./resource_pbix_report_only.pbix", reportPbixLocation, []pbixrewriter.PipelineFunc{
						pbixrewriter.SetDatasetIDPipelineFunc(datasetID),
					})
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "report_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test report PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
				}
				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, reportPbixLocationTfFriendly, reportPbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),
					testCheckReportExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),
					testCheckResourceAttrNotSet("powerbi_pbix.report_only", "dataset_id"),
				),
			},
			// Attempt to set parameters
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "report_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test report PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					parameter {
						name = "ParamOne"
						value = "NewParamValueOne"
					}
				}
				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, reportPbixLocationTfFriendly, reportPbixLocationTfFriendly),
				ExpectError: regexp.MustCompile("Unable to update parameters on a PBIX file that does not contain a dataset"),
			},
			// Attempt to set datasources
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "report_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test report PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					datasource {
						type = "OData"
						url = "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"
						original_url = "https://services.odata.org/V3/OData/OData.svc"
					}
				}
				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, reportPbixLocationTfFriendly, reportPbixLocationTfFriendly),
				ExpectError: regexp.MustCompile("Unable to update datasources on a PBIX file that does not contain a dataset"),
			},
			// deletes the resource
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test dataset PBIX"),

					testCheckReportDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),
					testCheckResourceRemoved("powerbi_pbix.dataset_only"),
					testCheckResourceRemoved("powerbi_pbix.report_only"),
				),
			},
		},
	})
}

func TestAccPBIX_rebind_dataset(t *testing.T) {
	datasetPbixLocation := TempFileName("dataset_", ".pbix")
	datasetPbixLocationTfFriendly := strings.ReplaceAll(datasetPbixLocation, "\\", "\\\\")

	reportPbixLocation := TempFileName("report_", ".pbix")
	reportPbixLocationTfFriendly := strings.ReplaceAll(reportPbixLocation, "\\", "\\\\")

	workspaceSuffix := acctest.RandString(6)
	var report_datasetID string
	var dataset_datasetID string
	var dataset2_datasetID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step create a dataset
			{
				PreConfig: func() {
					Copy("./resource_pbix_test_sample1.pbix", datasetPbixLocation)
					Copy("./resource_pbix_test_sample1.pbix", reportPbixLocation)
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "report_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test report PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					rebind_dataset_id = powerbi_pbix.dataset_only.dataset_id
				}
				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, reportPbixLocationTfFriendly, reportPbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.dataset_only", "dataset_id", &dataset_datasetID),
					testCheckReportDataset("powerbi_pbix.report_only", &dataset_datasetID),
				),
			},

			// second step rebind to a different report
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "dataset_only_2" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset 2 PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "report_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test report PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					rebind_dataset_id = powerbi_pbix.dataset_only_2.dataset_id
				}

				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, reportPbixLocationTfFriendly, reportPbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.dataset_only_2", "dataset_id", &dataset2_datasetID),
					testCheckReportDataset("powerbi_pbix.report_only", &dataset2_datasetID),
				),
			},

			// third step remove binding
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "dataset_only_2" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset 2 PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}

				resource "powerbi_pbix" "report_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test report PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
				}

				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly, reportPbixLocationTfFriendly, reportPbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.report_only", "dataset_id", &report_datasetID),
					testCheckReportDataset("powerbi_pbix.report_only", &report_datasetID),
				),
			},

			// fourth step delete the report
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "dataset_only" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test dataset PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
					skip_report = true
				}
				`, workspaceSuffix, datasetPbixLocationTfFriendly, datasetPbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetExistsInWorkspace("powerbi_workspace.test", "Acceptance Test dataset PBIX"),

					testCheckReportDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),

					testCheckResourceRemoved("powerbi_pbix.report_only"),
				),
			},

			// deletes all resources
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test dataset PBIX"),

					testCheckReportDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test report PBIX"),

					testCheckResourceRemoved("powerbi_pbix.dataset_only"),
					testCheckResourceRemoved("powerbi_pbix.report_only"),
				),
			},
		},
	})
}

func TestAccPBIX_parameters(t *testing.T) {
	var updatedTime time.Time
	var datasetID string
	var groupID string
	workspaceSuffix := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the pbix with parameters
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "./resource_pbix_test_sample1.pbix"
					source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
					parameter {
						name = "ParamOne"
						value = "NewParamValueOne"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.test", "dataset_id", &datasetID),
					set("powerbi_pbix.test", "workspace_id", &groupID),
					setUpdatedTime("powerbi_pbix.test", &updatedTime),
					testCheckParameter("powerbi_pbix.test", "ParamOne", "NewParamValueOne"),
				),
			},
			// identical resource definition with parameter state drift
			{
				PreConfig: func() {
					//update parameter outside of terraform to simulate drift
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateParametersInGroup(groupID, datasetID, powerbiapi.UpdateParametersInGroupRequest{
						UpdateDetails: []powerbiapi.UpdateParametersInGroupRequestItem{
							{
								Name:     "ParamOne",
								NewValue: "DriftedValue",
							},
							{
								Name:     "ParamTwo",
								NewValue: "DriftedValue",
							},
						},
					})
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "./resource_pbix_test_sample1.pbix"
					source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
					parameter {
						name = "ParamOne"
						value = "NewParamValueOne"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckUpdatedAt("powerbi_pbix.test", &updatedTime), //import should not be updated
					testCheckParameter("powerbi_pbix.test", "ParamOne", "NewParamValueOne"),
					testCheckParameter("powerbi_pbix.test", "ParamTwo", "DriftedValue"),
				),
			},
			// uploading new file should also update with parameters
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "./resource_pbix_test_sample2.pbix"
					source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
					parameter {
						name = "ParamOne"
						value = "NewParamValueOne"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckUpdatedAfter("powerbi_pbix.test", &updatedTime),                //import should be updated
					testCheckParameter("powerbi_pbix.test", "ParamOne", "NewParamValueOne"), //new value maintained
					testCheckParameter("powerbi_pbix.test", "ParamTwo", "ParamTwoValue"),
				),
			},
		},
	})
}

func TestAccPBIX_datasources(t *testing.T) {
	var updatedTime time.Time
	var datasetID string
	var groupID string
	workspaceSuffix := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the pbix with datasource change
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "./resource_pbix_test_sample1.pbix"
					source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
					datasource {
						type = "OData"
						url = "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"
						original_url = "https://services.odata.org/V3/OData/OData.svc"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.test", "dataset_id", &datasetID),
					set("powerbi_pbix.test", "workspace_id", &groupID),
					setUpdatedTime("powerbi_pbix.test", &updatedTime),
					testCheckURLDatasource("powerbi_pbix.test", "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"),
				),
			},
			// apply same config with drift
			{
				PreConfig: func() {
					//update datasource outside of terraform to simulate drift
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateDatasourcesInGroup(groupID, datasetID, powerbiapi.UpdateDatasourcesInGroupRequest{
						UpdateDetails: []powerbiapi.UpdateDatasourcesInGroupRequestItem{
							{
								ConnectionDetails: powerbiapi.UpdateDatasourcesInGroupRequestItemConnectionDetails{
									URL: emptyStringToNil("https://google.com"),
								},
								DatasourceSelector: powerbiapi.UpdateDatasourcesInGroupRequestItemDatasourceSelector{
									DatasourceType: "OData",
									ConnectionDetails: powerbiapi.UpdateDatasourcesInGroupRequestItemConnectionDetails{
										URL: emptyStringToNil("https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"),
									},
								},
							},
						},
					})
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "./resource_pbix_test_sample1.pbix"
					source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
					datasource {
						type = "OData"
						url = "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"
						original_url = "https://services.odata.org/V3/OData/OData.svc"
					}
				}
				`, workspaceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckUpdatedAfter("powerbi_pbix.test", &updatedTime), //import should be updated
					testCheckURLDatasource("powerbi_pbix.test", "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"),
				),
			},
		},
	})
}

// TempFileName generates a temporary filename for use in testing or whatever
func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func testCheckResourceRemoved(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if res, ok := s.RootModule().Resources[resourceName]; ok {
			return fmt.Errorf("Expected resource %v to be deleted. Resource still exists %v", resourceName, res)
		}
		return nil
	}
}

func getResourceID(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource not found: %s", resourceName)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("resource id not set")
	}
	return rs.Primary.ID, nil
}

func getResourceProperty(s *terraform.State, resourceName string, property string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource not found: %s", resourceName)
	}

	propVal, ok := rs.Primary.Attributes[property]
	if !ok {
		return "", fmt.Errorf("resource property %s not set", property)
	}
	return propVal, nil
}

func set(resourceName string, configName string, outID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		*outID = rs.Primary.Attributes[configName]
		return nil
	}
}

func setUpdatedTime(pbixResourceName string, outUpdatedTime *time.Time) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pbixID, err := getResourceID(s, pbixResourceName)
		if err != nil {
			return err
		}

		groupID, err := getResourceProperty(s, pbixResourceName, "workspace_id")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		im, err := client.GetImportInGroup(groupID, pbixID)
		if err != nil {
			return err
		}

		*outUpdatedTime = im.UpdatedDateTime

		return nil
	}
}

func testCheckDatasetExistsInWorkspace(workspaceResourceName string, expectedDatasetName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupID, err := getResourceID(s, workspaceResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		datasets, err := client.GetDatasetsInGroup(groupID)
		if err != nil {
			return err
		}

		var datasetNames []string
		for _, dataset := range datasets.Value {
			if dataset.Name == expectedDatasetName {
				return nil
			}
			datasetNames = append(datasetNames, dataset.Name)
		}
		return fmt.Errorf("Expecting datasets %v in workspace %v. Found datasets %v", expectedDatasetName, groupID, datasetNames)
	}
}

func testCheckDatasetDoesNotExistsInWorkspace(workspaceResourceName string, expectedDatasetName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupID, err := getResourceID(s, workspaceResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		datasets, err := client.GetDatasetsInGroup(groupID)
		if err != nil {
			return err
		}

		for _, dataset := range datasets.Value {
			if dataset.Name == expectedDatasetName {
				return fmt.Errorf("Expecting no datasets in workspace %v. Found datasets %v", groupID, dataset.Name)
			}
		}
		return nil
	}
}

func testCheckReportExistsInWorkspace(workspaceResourceName string, expectedReportName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupID, err := getResourceID(s, workspaceResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		reports, err := client.GetReportsInGroup(groupID)
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
		return fmt.Errorf("Expecting reports %v in workspace %s. Found reports %v", expectedReportName, groupID, reportNames)
	}
}

func testCheckReportDoesNotExistsInWorkspace(workspaceResourceName string, expectedReportName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupID, err := getResourceID(s, workspaceResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		reports, err := client.GetReportsInGroup(groupID)

		if err != nil {
			return err
		}

		for _, report := range reports.Value {
			if report.Name == expectedReportName {
				return fmt.Errorf("Expecting report to not exist in workspace %v. Found report %v", groupID, report.Name)
			}
		}
		return nil
	}
}

func testCheckReportDataset(pbixResourceName string, expectedDatasetId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		reportID, err := getResourceProperty(s, pbixResourceName, "report_id")
		if err != nil {
			return err
		}

		groupID, err := getResourceProperty(s, pbixResourceName, "workspace_id")
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		report, err := client.GetReportInGroup(groupID, reportID)
		if err != nil {
			return err
		}
		if report.DatasetID != *expectedDatasetId {
			return fmt.Errorf("Expecting report %v to have to dataset %v. Found report to have dataset %v", reportID, *expectedDatasetId, report.DatasetID)
		}

		return nil
	}
}

func testCheckResourceAttrNotSet(name string, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		propVal, ok := rs.Primary.Attributes[key]
		if ok {
			return fmt.Errorf("Expected property %s to not be set. Found property %s with value %s", key, key, propVal)
		}

		return nil
	}
}

func testCheckUpdatedAfter(pbixResourceName string, updatedAfter *time.Time) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pbixID, err := getResourceID(s, pbixResourceName)
		if err != nil {
			return err
		}
		groupID, err := getResourceProperty(s, pbixResourceName, "workspace_id")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		im, err := client.GetImportInGroup(groupID, pbixID)
		if err != nil {
			return err
		}

		if updatedAfter != nil && (im.UpdatedDateTime.Before(*updatedAfter) || im.UpdatedDateTime.Equal(*updatedAfter)) {
			return fmt.Errorf("Expected to find import %v updated after %v. Found import updated at %v", pbixID, updatedAfter, im.UpdatedDateTime)
		}

		return nil

	}
}

func testCheckUpdatedAt(pbixResourceName string, updatedAt *time.Time) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pbixID, err := getResourceID(s, pbixResourceName)
		if err != nil {
			return err
		}
		groupID, err := getResourceProperty(s, pbixResourceName, "workspace_id")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		im, err := client.GetImportInGroup(groupID, pbixID)
		if err != nil {
			return err
		}

		if !im.UpdatedDateTime.Equal(*updatedAt) {
			return fmt.Errorf("Expected to find import %v updated at %v. Found import updated at %v", pbixID, updatedAt, im.UpdatedDateTime)
		}

		return nil

	}
}

func testCheckParameter(pbixResourceName string, expectedParameterName string, expectedParameterValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		datasetID, err := getResourceProperty(s, pbixResourceName, "dataset_id")
		if err != nil {
			return err
		}

		groupID, err := getResourceProperty(s, pbixResourceName, "workspace_id")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		params, err := client.GetParametersInGroup(groupID, datasetID)
		if err != nil {
			return err
		}

		var parameterNames []string
		for _, param := range params.Value {
			parameterNames = append(parameterNames, param.Name)
			if param.Name == expectedParameterName {
				if param.CurrentValue != expectedParameterValue {
					return fmt.Errorf("Expecting parameter %v to have a value of %s. Found value of %v", expectedParameterName, expectedParameterValue, param.CurrentValue)
				}
				return nil
			}
		}

		return fmt.Errorf("Expecting parameter with name %s to exist. Only the following parameters %s were found", expectedParameterName, parameterNames)
	}
}

func testCheckURLDatasource(pbixResourceName string, expectedValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		datasetID, err := getResourceProperty(s, pbixResourceName, "dataset_id")
		if err != nil {
			return err
		}
		groupID, err := getResourceProperty(s, pbixResourceName, "workspace_id")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		datasources, err := client.GetDatasourcesInGroup(groupID, datasetID)
		if err != nil {
			return err
		}

		var urlValues []string
		for _, datasource := range datasources.Value {
			if *datasource.ConnectionDetails.URL == expectedValue {
				return nil
			} else if datasource.ConnectionDetails.URL != nil {
				urlValues = append(urlValues, *datasource.ConnectionDetails.URL)
			}
		}

		return fmt.Errorf("Expecting datasource with field url value %s to exist. Only the urls %v were found in the datasources", expectedValue, urlValues)
	}
}
