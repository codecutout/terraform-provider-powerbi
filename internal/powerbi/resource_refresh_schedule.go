package powerbi

import (
	"fmt"
	"regexp"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceRefreshSchedule represents a Power BI refresh schedule
func ResourceRefreshSchedule() *schema.Resource {
	return &schema.Resource{
		Create: createRefreshSchedule,
		Read:   readRefreshSchedule,
		Update: updateRefreshSchedule,
		Delete: deleteRefreshSchedule,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "Workspace ID in which the dataset was deployed.",
				Required:    true,
				ForceNew:    true,
			},
			"dataset_id": {
				Type:        schema.TypeString,
				Description: "The ID for the dataset that was deployed as part of the PBIX.",
				Required:    true,
				ForceNew:    true,
			},
			"days": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of days of the week when the schedule should refresh.",
				Required:    true,
			},
			"times": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of times on the day the schedule should refresh. Times should be in the format HH:00 or HH:30 i.e. Hour should be two digits and minutes must either be on the full or half hour.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Determines if the scheduled refresh is enabled.",
				Optional:    true,
				Default:     true,
			},
			"local_time_zone_id": {
				Type:        schema.TypeString,
				Description: "The name of the timezone to use. See Name of Time Zone column in [Microsoft Time Zone Index Values](https://support.microsoft.com/en-gb/help/973627/microsoft-time-zone-index-values).",
				Optional:    true,
				Default:     "UTC",
			},
			"notify_option": {
				Type:        schema.TypeString,
				Description: "The notification option when a scheduled refresh fails. Should be either `MailOnFailure` or `NoNotification`",
				Optional:    true,
				Default:     "NoNotification",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					stringVal := val.(string)
					reg := regexp.MustCompile("^(MailOnFailure|NoNotification)$")
					if !reg.MatchString(stringVal) {
						errs = append(errs, fmt.Errorf("Expected argument 'notify_option' to be either 'MailOnFailure' or 'NoNotification'. Found '%v'", stringVal))
					}
					return warns, errs
				},
			},
		},
	}
}

func getDatasetID(d *schema.ResourceData, meta interface{}) (string, error) {
	datasetID := d.Get("dataset_id").(string)
	if datasetID == "" {
		datasetID = d.Id()
	}
	if datasetID == "" {
		return "", fmt.Errorf("Unable to determine dataset ID. Ensure dataset_id is set")
	}
	return datasetID, nil
}

func getGroupID(d *schema.ResourceData, meta interface{}) (string, error) {
	groupID := d.Get("workspace_id").(string)

	if groupID == "" {
		return "", fmt.Errorf("Unable to determine workspace ID. Ensure workspace_id is set")
	}

	return groupID, nil
}

func validateConfig(d *schema.ResourceData, meta interface{}) error {
	// schema validate functions do not yet support lists and maps
	// creating own makeshift validation to check days and times

	err := validateConfigDays(d, meta)
	if err != nil {
		return err
	}
	err = validateConfigTimes(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func validateConfigDays(d *schema.ResourceData, meta interface{}) error {
	reg := regexp.MustCompile("^(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)$")
	days := convertToStringSlice(d.Get("days").([]interface{}))
	for _, day := range days {
		if !reg.MatchString(day) {
			return fmt.Errorf("config is invalid: Expected argument 'days' to be either 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday' or 'Sunday'. Found '%v'", day)
		}
	}
	return nil
}

func validateConfigTimes(d *schema.ResourceData, meta interface{}) error {
	reg := regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]):(00|30)$")
	times := convertToStringSlice(d.Get("times").([]interface{}))
	for _, time := range times {
		if !reg.MatchString(time) {
			return fmt.Errorf("config is invalid: Expected argument 'times' to be in the format 'HH:00' or 'HH:30'. Hours must be two digits and must be on the hour or half hour. Found time '%v'", time)
		}
	}
	return nil
}

func createRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	err := validateConfig(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*powerbiapi.Client)

	enabled := nilIfFalse(d.Get("enabled").(bool))
	datasetID, err := getDatasetID(d, meta)
	if err != nil {
		return err
	}

	groupID, err := getGroupID(d, meta)
	if err != nil {
		return err
	}

	err = client.UpdateRefreshScheduleInGroup(groupID, datasetID, powerbiapi.UpdateRefreshScheduleInGroupRequest{
		Value: powerbiapi.UpdateRefreshScheduleInGroupRequestValue{
			Enabled:         convertBoolToPointer(true), // API doesnt allow updating if disabled
			Days:            convertStringSliceToPointer(convertToStringSlice(d.Get("days").([]interface{}))),
			Times:           convertStringSliceToPointer(convertToStringSlice(d.Get("times").([]interface{}))),
			LocalTimeZoneID: convertStringToPointer(d.Get("local_time_zone_id").(string)),
			NotifyOption:    convertStringToPointer(d.Get("notify_option").(string)),
		},
	})
	if err != nil {
		return err
	}

	// Set the disabled flag to be the correct value
	if enabled == nil {
		err := client.UpdateRefreshScheduleInGroup(groupID, datasetID, powerbiapi.UpdateRefreshScheduleInGroupRequest{
			Value: powerbiapi.UpdateRefreshScheduleInGroupRequestValue{
				Enabled: convertBoolToPointer(false),
			},
		})
		if err != nil {
			return err
		}
	}

	return readRefreshSchedule(d, meta)
}

func readRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	datasetID, err := getDatasetID(d, meta)
	if err != nil {
		return err
	}
	groupID, err := getGroupID(d, meta)
	if err != nil {
		return err
	}

	refreshSchedule, err := client.GetRefreshScheduleInGroup(groupID, datasetID)
	if isHTTP404Error(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.SetId(datasetID)
	d.Set("dataset_id", datasetID)
	d.Set("workspace_id", groupID)
	d.Set("enabled", refreshSchedule.Enabled)
	d.Set("days", refreshSchedule.Days)
	d.Set("times", refreshSchedule.Times)
	d.Set("local_time_zone_id", refreshSchedule.LocalTimeZoneID)
	d.Set("notify_option", refreshSchedule.NotifyOption)

	return nil
}

func updateRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	err := validateConfig(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*powerbiapi.Client)

	requestVal := powerbiapi.UpdateRefreshScheduleInGroupRequestValue{}
	updateRequired := false
	disableRequired := false

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		if enabled {
			requestVal.Enabled = convertBoolToPointer(d.Get("enabled").(bool))
			updateRequired = true
		} else {
			disableRequired = true
		}
	}
	if d.HasChange("days") {
		requestVal.Days = convertStringSliceToPointer(convertToStringSlice(d.Get("days").([]interface{})))
		updateRequired = true
	}
	if d.HasChange("times") {
		requestVal.Times = convertStringSliceToPointer(convertToStringSlice(d.Get("times").([]interface{})))
		updateRequired = true
	}
	if d.HasChange("local_time_zone_id") {
		requestVal.LocalTimeZoneID = convertStringToPointer(d.Get("local_time_zone_id").(string))
		updateRequired = true
	}
	if d.HasChange("notify_option") {
		requestVal.NotifyOption = convertStringToPointer(d.Get("notify_option").(string))
		updateRequired = true
	}

	datasetID, err := getDatasetID(d, meta)
	if err != nil {
		return err
	}
	groupID, err := getGroupID(d, meta)
	if err != nil {
		return err
	}

	if updateRequired {
		err := client.UpdateRefreshScheduleInGroup(groupID, datasetID, powerbiapi.UpdateRefreshScheduleInGroupRequest{
			Value: requestVal,
		})
		if err != nil {
			return err
		}
	}

	// disabling has to be in a seperate step as api does not allow updates and disable in same request
	if disableRequired {
		err := client.UpdateRefreshScheduleInGroup(groupID, datasetID, powerbiapi.UpdateRefreshScheduleInGroupRequest{
			Value: powerbiapi.UpdateRefreshScheduleInGroupRequestValue{Enabled: convertBoolToPointer(false)},
		})
		if err != nil {
			return err
		}
	}

	return readRefreshSchedule(d, meta)
}

func deleteRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	// You dont delete refresh schedules, so we will disable it
	datasetID, err := getDatasetID(d, meta)
	if err != nil {
		return err
	}
	groupID, err := getGroupID(d, meta)
	if err != nil {
		return err
	}

	return client.UpdateRefreshScheduleInGroup(groupID, datasetID, powerbiapi.UpdateRefreshScheduleInGroupRequest{
		Value: powerbiapi.UpdateRefreshScheduleInGroupRequestValue{
			Enabled: convertBoolToPointer(false),
		},
	})
}
