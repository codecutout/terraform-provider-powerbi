package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/powerbi/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceRefreshSchedule represents a Power BI refresh schedule
func ResourceRefreshSchedule() *schema.Resource {
	return &schema.Resource{
		Create: createRefreshSchedule,
		Read:   readRefreshSchedule,
		Update: updateRefreshSchedule,
		Delete: deleteRefreshSchedule,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
				Description: "The list of times on the day the schedule should refresh. Should only be in half hour increments.",
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
			},
		},
	}
}

func convertStringToPointer(s string) *string {
	return &s
}

func convertBoolToPointer(b bool) *bool {
	return &b
}

func convertStringSliceToPointer(ss []string) *[]string {
	return &ss
}

func convertToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i := range interfaceSlice {
		stringSlice[i] = interfaceSlice[i].(string)
	}
	return stringSlice
}

func nilIfFalse(b bool) *bool {
	if !b {
		return nil
	}
	return &b

}

func createRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	enabled := nilIfFalse(d.Get("enabled").(bool))
	datasetID := d.Get("dataset_id").(string)
	err := client.UpdateRefreshSchedule(datasetID, api.UpdateRefreshScheduleRequest{
		Value: api.UpdateRefreshScheduleRequestValue{
			Enabled:         enabled,
			Days:            convertStringSliceToPointer(convertToStringSlice(d.Get("days").([]interface{}))),
			Times:           convertStringSliceToPointer(convertToStringSlice(d.Get("times").([]interface{}))),
			LocalTimeZoneID: convertStringToPointer(d.Get("local_time_zone_id").(string)),
			NotifyOption:    convertStringToPointer(d.Get("notify_option").(string)),
		},
	})
	if err != nil {
		return err
	}

	// API does not allow disabling while changing other properties
	if enabled == nil {
		err := client.UpdateRefreshSchedule(datasetID, api.UpdateRefreshScheduleRequest{
			Value: api.UpdateRefreshScheduleRequestValue{
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
	client := meta.(*api.Client)

	datasetID := d.Get("dataset_id").(string)
	refreshSchedule, err := client.GetRefreshSchedule(datasetID)
	if err != nil {
		return err
	}

	d.SetId(datasetID)
	d.Set("enabled", refreshSchedule.Enabled)
	d.Set("days", refreshSchedule.Days)
	d.Set("times", refreshSchedule.Times)
	d.Set("local_time_zone_id", refreshSchedule.LocalTimeZoneID)
	d.Set("notify_option", refreshSchedule.NotifyOption)

	return nil
}

func updateRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	requestVal := api.UpdateRefreshScheduleRequestValue{}
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

	datasetID := d.Get("dataset_id").(string)
	if updateRequired {
		err := client.UpdateRefreshSchedule(datasetID, api.UpdateRefreshScheduleRequest{
			Value: requestVal,
		})
		if err != nil {
			return err
		}
	}

	// disabling has to be in a seperate step as api does not allow updates and disable in same request
	if disableRequired {
		err := client.UpdateRefreshSchedule(datasetID, api.UpdateRefreshScheduleRequest{
			Value: api.UpdateRefreshScheduleRequestValue{Enabled: convertBoolToPointer(false)},
		})
		if err != nil {
			return err
		}
	}

	return readRefreshSchedule(d, meta)
}

func deleteRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	// You dont delete refresh schedules, so we will disable it
	datasetID := d.Get("dataset_id").(string)
	err := client.UpdateRefreshSchedule(datasetID, api.UpdateRefreshScheduleRequest{
		Value: api.UpdateRefreshScheduleRequestValue{
			Enabled: convertBoolToPointer(false),
		},
	})

	if err != nil {
		// we do not care about 404. Indictes the resource is already deleted
		httpErr, isHTTPErr := err.(api.HTTPUnsuccessfulError)
		if isHTTPErr && httpErr.Response.StatusCode == 404 {
			return nil
		}
	}

	return err
}
