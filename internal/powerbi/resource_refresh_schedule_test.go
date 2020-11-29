package powerbi

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRefreshSchedule_basic(t *testing.T) {
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
					workspace_id = "${powerbi_workspace.test.id}"
					name = "Acceptance Test PBIX"
					source = "./resource_pbix_test_sample1.pbix"
					source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
				}

				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "${powerbi_pbix.test.dataset_id}"
					enabled = true
					days = ["Monday", "Wednesday", "Friday"]
					times = ["09:00", "17:30"]
					local_time_zone_id = "Pacific Standard Time"
					notify_option = "MailOnFailure"
				}

				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "days.0", "Monday"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "days.1", "Wednesday"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "days.2", "Friday"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "times.0", "09:00"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "times.1", "17:30"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "local_time_zone_id", "Pacific Standard Time"),
					resource.TestCheckResourceAttr("powerbi_refresh_schedule.test", "notify_option", "MailOnFailure"),
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleResponse{
						Enabled:         true,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "MailOnFailure",
					}),
				),
			},
			// second step updates the resource
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
				}

				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "${powerbi_pbix.test.dataset_id}"
					enabled = true
					days = ["Tuesday", "Thursday"] # days changed
					times = ["09:00", "17:30"]
					local_time_zone_id = "UTC" # time zone changed
					notify_option = "MailOnFailure"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleResponse{
						Enabled:         true,
						Days:            []string{"Tuesday", "Thursday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "UTC",
						NotifyOption:    "MailOnFailure",
					}),
				),
			},
			// third step updates the rest of the resource
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
				}

				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "${powerbi_pbix.test.dataset_id}"
					enabled = false # enabled changed
					days = ["Tuesday", "Thursday"] 
					times = ["09:00"] # times changed
					local_time_zone_id = "UTC"
					notify_option = "NoNotification" # notify option changed
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleResponse{
						Enabled:         false,
						Days:            []string{"Tuesday", "Thursday"},
						Times:           []string{"09:00"},
						LocalTimeZoneID: "UTC",
						NotifyOption:    "NoNotification",
					}),
				),
			},
			// final step checks importing the current state we reached in the step above
			{
				ResourceName:      "powerbi_refresh_schedule.test",
				ImportState:       true,
				ImportStateVerify: true,
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
					times = []
					days = []
					notify_option = "not-an-option"
				}

				`,
				ExpectError: regexp.MustCompile("config is invalid:.*notify_option.*"),
			},
			{
				Config: `
				resource "powerbi_refresh_schedule" "test" {
					dataset_id = "validation-should-fail-before-using-this"
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

	config := `
	resource "powerbi_workspace" "test" {
		name = "Acceptance Test Workspace"
	}

	resource "powerbi_pbix" "test" {
		workspace_id = "${powerbi_workspace.test.id}"
		name = "Acceptance Test PBIX"
		source = "./resource_pbix_test_sample1.pbix"
		source_hash = "${filemd5("./resource_pbix_test_sample1.pbix")}"
	}

	resource "powerbi_refresh_schedule" "test" {
		dataset_id = "${powerbi_pbix.test.dataset_id}"
		enabled = false
		days = ["Monday", "Wednesday", "Friday"]
		times = ["09:00", "17:30"]
		local_time_zone_id = "Pacific Standard Time"
		notify_option = "MailOnFailure"
	}

	`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first creates the resource
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_refresh_schedule.test", "dataset_id", &datasetID),
				),
			},
			// second step skew the resource and checks it gets reupdates it
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateRefreshSchedule(datasetID, powerbiapi.UpdateRefreshScheduleRequest{
						Value: powerbiapi.UpdateRefreshScheduleRequestValue{
							LocalTimeZoneID: convertStringToPointer("UTC"),
						},
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleResponse{
						Enabled:         false,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "MailOnFailure",
					}),
				),
			},
			// third step deletes dataset
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.DeleteDataset(datasetID)
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", powerbiapi.GetRefreshScheduleResponse{
						Enabled:         false,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "MailOnFailure",
					}),
				),
			},
		},
	})
}

func testCheckRefreshSchedule(scheduleRefreshResourceName string, expectedRefreshSchedule powerbiapi.GetRefreshScheduleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		refreshScheduleID, err := getID(s, scheduleRefreshResourceName)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*powerbiapi.Client)
		actualRefreshSchedule, err := client.GetRefreshSchedule(refreshScheduleID)

		if err != nil {
			return err
		}

		if !reflect.DeepEqual(expectedRefreshSchedule, *actualRefreshSchedule) {
			return fmt.Errorf("Expected refresh schedule %v. Found refresh schedule %v", expectedRefreshSchedule, actualRefreshSchedule)
		}

		return nil
	}
}
