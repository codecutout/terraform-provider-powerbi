package powerbi

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRefreshSchedule_basic(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
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
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_pbix" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "%5s"
					source_hash = "${filemd5("%s")}"
				}

				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "${powerbi_pbix.test.dataset_id}"
					workspace_id = "${powerbi_pbix.test.workspace_id}"
					enabled = true
					days = ["Monday", "Wednesday", "Friday"]
					times = ["09:00", "17:30"]
					local_time_zone_id = "Pacific Standard Time"
					notify_option = "NoNotification"
				}

				`, workspaceSuffix, pbixLocationTfFriendly, pbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "days.0", "Monday"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "days.1", "Wednesday"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "days.2", "Friday"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "times.0", "09:00"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "times.1", "17:30"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "local_time_zone_id", "Pacific Standard Time"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "notify_option", "NoNotification"),
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleInGroupResponse{
						Enabled:         true,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "NoNotification",
					}),
				),
			},
			// second step updates the resource
			{
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

				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "${powerbi_pbix.test.dataset_id}"
					workspace_id = "${powerbi_pbix.test.workspace_id}"
					enabled = true
					days = ["Tuesday", "Thursday"] # days changed
					times = ["09:00", "17:30"]
					local_time_zone_id = "UTC" # time zone changed
					notify_option = "NoNotification"
				}
				`, workspaceSuffix, pbixLocationTfFriendly, pbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleInGroupResponse{
						Enabled:         true,
						Days:            []string{"Tuesday", "Thursday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "UTC",
						NotifyOption:    "NoNotification",
					}),
				),
			},
			// third step updates the rest of the resource
			{
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

				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "${powerbi_pbix.test.dataset_id}"
					workspace_id = "${powerbi_pbix.test.workspace_id}"
					enabled = false # enabled changed
					days = ["Tuesday", "Friday"] # days changed 
					times = ["09:00"] # times changed
					local_time_zone_id = "UTC"
					notify_option = "NoNotification"
				}
				`, workspaceSuffix, pbixLocationTfFriendly, pbixLocationTfFriendly),
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleInGroupResponse{
						Enabled:         false,
						Days:            []string{"Tuesday", "Friday"},
						Times:           []string{"09:00"},
						LocalTimeZoneID: "UTC",
						NotifyOption:    "NoNotification",
					}),
				),
			},
		},
	})
}

func TestAccRefreshSchedule_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "validation-should-fail-before-using-this"
					workspace_id = "validation-should-fail-before-using-this"
					times = []
					days = []
					notify_option = "not-an-option"
				}

				`,
				ExpectError: regexp.MustCompile(".*notify_option.*not-an-option"),
			},
			{
				Config: `
				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "validation-should-fail-before-using-this"
					workspace_id = "validation-should-fail-before-using-this"
					times = []
					days = ["Monday", "Badday", "Wednesday"]
				}

				`,
				ExpectError: regexp.MustCompile("config is invalid:.*days.*"),
			},
			{
				Config: `
				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "validation-should-fail-before-using-this"
					workspace_id = "validation-should-fail-before-using-this"
					times = ["9:30"]
					days = []
				}

				`,
				ExpectError: regexp.MustCompile("config is invalid:.*times.*"),
			},
			{
				Config: `
				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "validation-should-fail-before-using-this"
					workspace_id = "validation-should-fail-before-using-this"
					times = ["09:45"]
					days = []
				}

				`,
				ExpectError: regexp.MustCompile("config is invalid:.*times.*"),
			},
		},
	})
}

func TestAccRefreshSchedule_skew(t *testing.T) {
	var datasetID string
	var groupID string
	workspaceSuffix := acctest.RandString(6)
	pbixLocation := TempFileName("", ".pbix")
	pbixLocationTfFriendly := strings.ReplaceAll(pbixLocation, "\\", "\\\\")

	config := fmt.Sprintf(`
	resource "powerbi_workspace" "test" {
		name = "Acceptance Test Workspace %s"
	}

	resource "powerbi_pbix" "test" {
		workspace_id = "${powerbi_workspace.test.id}"
		name = "Acceptance Test PBIX"
		source = "%s"
		source_hash = "${filemd5("%s")}"
	}

	resource "powerbi_refresh_schedule" "test" {
		dataset_id = "${powerbi_pbix.test.dataset_id}"
		workspace_id = "${powerbi_pbix.test.workspace_id}"
		enabled = false
		days = ["Monday", "Wednesday", "Friday"]
		times = ["09:00", "17:30"]
		local_time_zone_id = "Pacific Standard Time"
		notify_option = "NoNotification"
	}

	`, workspaceSuffix, pbixLocationTfFriendly, pbixLocationTfFriendly)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first creates the resource
			{
				PreConfig: func() {
					Copy("./resource_pbix_test_sample1.pbix", pbixLocation)
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_refresh_schedule.test", "dataset_id", &datasetID),
					set("powerbi_refresh_schedule.test", "worksapce_id", &groupID),
				),
			},
			// second step skew the resource and checks it gets reupdates it
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateRefreshScheduleInGroup(groupID, datasetID, powerbiapi.UpdateRefreshScheduleInGroupRequest{
						Value: powerbiapi.UpdateRefreshScheduleInGroupRequestValue{
							LocalTimeZoneID: convertStringToPointer("UTC"),
						},
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleInGroupResponse{
						Enabled:         false,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "NoNotification",
					}),
				),
			},
			// third step deletes dataset
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.DeleteDatasetInGroup(groupID, datasetID)
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleInGroupResponse{
						Enabled:         false,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "NoNotification",
					}),
				),
			},
		},
	})
}

func testCheckRefreshSchedule(scheduleRefreshResourceName string, expectedRefreshSchedule powerbiapi.GetRefreshScheduleInGroupResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		datasetID, err := getResourceProperty(s, scheduleRefreshResourceName, "dataset_id")
		if err != nil {
			return err
		}
		groupID, err := getResourceProperty(s, scheduleRefreshResourceName, "workspace_id")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		actualRefreshSchedule, err := client.GetRefreshScheduleInGroup(groupID, datasetID)

		if err != nil {
			return err
		}

		if !reflect.DeepEqual(expectedRefreshSchedule, *actualRefreshSchedule) {
			return fmt.Errorf("Expected refresh schedule %v. Found refresh schedule %v", expectedRefreshSchedule, actualRefreshSchedule)
		}

		return nil
	}
}
