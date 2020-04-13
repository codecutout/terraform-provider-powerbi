package powerbi

import (
	"fmt"
	"github.com/codecutout/terraform-provider-powerbi/powerbi/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"reflect"
	"testing"
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
					testCheckRefreshSchedule("powerbi_refresh_schedule.test", api.GetRefreshScheduleResponse{
						Enabled:         true,
						Days:            []string{"Monday", "Wednesday", "Friday"},
						Times:           []string{"09:00", "17:30"},
						LocalTimeZoneID: "Pacific Standard Time",
						NotifyOption:    "MailOnFailure",
					}),
				),
			},
			// second step updates it with a new title
			// {
			// 	Config: `
			// 	resource "powerbi_workspace" "test" {
			// 		name = "Acceptance Test Workspace - Updated"
			// 	}
			// 	`,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testCheckWorkspaceExistsWithName("powerbi_workspace.test", "Acceptance Test Workspace - Updated"),
			// 		resource.TestCheckResourceAttrSet("powerbi_workspace.test", "id"),
			// 		resource.TestCheckResourceAttr("powerbi_workspace.test", "name", "Acceptance Test Workspace - Updated"),
			// 	),
			// },
			// // final step checks importing the current state we reached in the step above
			// {
			// 	ResourceName:      "powerbi_workspace.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testCheckRefreshSchedule(scheduleRefreshResourceName string, expectedRefreshSchedule api.GetRefreshScheduleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		refreshScheduleID, err := getID(s, scheduleRefreshResourceName)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*api.Client)
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
