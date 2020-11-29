package powerbi

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPBIX_basic(t *testing.T) {
	var updatedTime time.Time
	pbixLocation := TempFileName("", ".pbix")
	pbixLocationTfFriendly := strings.ReplaceAll(pbixLocation, "\\", "\\\\")

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
					name = "Acceptance Test Workspace"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
				}
				`, pbixLocationTfFriendly, pbixLocationTfFriendly),
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
			// update wtih different pbix same path
			{
				PreConfig: func() {
					Copy("./resource_pbix_test_sample2.pbix", pbixLocation)
				},
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "%s"
					source_hash = "${filemd5("%s")}"
				}
				`, pbixLocationTfFriendly, pbixLocationTfFriendly),
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
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckDatasetDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckReportDoesNotExistsInWorkspace("powerbi_workspace.test", "Acceptance Test PBIX"),
					testCheckResourceRemoved("powerbi_pbix.test"),
				),
			},
		},
	})
}

func TestAccPBIX_parameters(t *testing.T) {
	var updatedTime time.Time
	var datasetID string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the pbix with parameters
			{
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
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
				`,
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.test", "dataset_id", &datasetID),
					setUpdatedTime("powerbi_pbix.test", &updatedTime),
					testCheckParameter("powerbi_pbix.test", "ParamOne", "NewParamValueOne"),
				),
			},
			// identical resource definition with parameter state drift
			{
				PreConfig: func() {
					//update paramter outside of terraform to simulate drift
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateParameters(datasetID, powerbiapi.UpdateParametersRequest{
						UpdateDetails: []powerbiapi.UpdateParametersRequestItem{
							powerbiapi.UpdateParametersRequestItem{
								Name:     "ParamOne",
								NewValue: "DriftedValue",
							},
							powerbiapi.UpdateParametersRequestItem{
								Name:     "ParamTwo",
								NewValue: "DriftedValue",
							},
						},
					})
				},
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
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
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckUpdatedAt("powerbi_pbix.test", &updatedTime), //import should not be updated
					testCheckParameter("powerbi_pbix.test", "ParamOne", "NewParamValueOne"),
					testCheckParameter("powerbi_pbix.test", "ParamTwo", "DriftedValue"),
				),
			},
			// uploading new file should also update with parameters
			{
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
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
				`,
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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the pbix with datasource change
			{
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
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
				`,
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_pbix.test", "dataset_id", &datasetID),
					setUpdatedTime("powerbi_pbix.test", &updatedTime),
					testCheckURLDatasource("powerbi_pbix.test", "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"),
				),
			},
			// apply same config with drift
			{
				PreConfig: func() {
					//update datasource outside of terraform to simulate drift
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateDatasources(datasetID, powerbiapi.UpdateDatasourcesRequest{
						UpdateDetails: []powerbiapi.UpdateDatasourcesRequestItem{
							powerbiapi.UpdateDatasourcesRequestItem{
								ConnectionDetails: powerbiapi.UpdateDatasourcesRequestItemConnectionDetails{
									URL: emptyStringToNil("https://google.com"),
								},
								DatasourceSelector: powerbiapi.UpdateDatasourcesRequestItemDatasourceSelector{
									DatasourceType: "OData",
									ConnectionDetails: powerbiapi.UpdateDatasourcesRequestItemConnectionDetails{
										URL: emptyStringToNil("https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"),
									},
								},
							},
						},
					})
				},
				Config: `
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace"
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
				`,
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

func getID(s *terraform.State, workspaceResourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[workspaceResourceName]
	if !ok {
		return "", fmt.Errorf("resource not found: %s", workspaceResourceName)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("resource id not set")
	}
	return rs.Primary.ID, nil
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
		pbixID, err := getID(s, pbixResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		im, err := client.GetImport(pbixID)
		if err != nil {
			return err
		}

		*outUpdatedTime = im.UpdatedDateTime

		return nil
	}
}

func testCheckDatasetExistsInWorkspace(workspaceResourceName string, expectedDatasetName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupID, err := getID(s, workspaceResourceName)
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
		groupID, err := getID(s, workspaceResourceName)
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
		groupID, err := getID(s, workspaceResourceName)
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
		groupID, err := getID(s, workspaceResourceName)
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

func testCheckUpdatedAfter(pbixResourceName string, updatedAfter *time.Time) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pbixID, err := getID(s, pbixResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		im, err := client.GetImport(pbixID)
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
		pbixID, err := getID(s, pbixResourceName)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*powerbiapi.Client)
		im, err := client.GetImport(pbixID)
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

		rs, ok := s.RootModule().Resources[pbixResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", pbixResourceName)
		}

		datasetID, ok := rs.Primary.Attributes["dataset_id"]
		if !ok {
			return fmt.Errorf("unable to find dataset_id on resource %s", pbixResourceName)
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		params, err := client.GetParameters(datasetID)
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

		rs, ok := s.RootModule().Resources[pbixResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", pbixResourceName)
		}

		datasetID, ok := rs.Primary.Attributes["dataset_id"]
		if !ok {
			return fmt.Errorf("unable to find dataset_id on resource %s", pbixResourceName)
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		datasources, err := client.GetDatasources(datasetID)
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
