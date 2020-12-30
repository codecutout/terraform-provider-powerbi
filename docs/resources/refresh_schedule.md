# Refresh Schedule Resource
`powerbi_refresh_schedule` represents a dataset's refresh schedule


## Example Usage
```hcl
resource "powerbi_refresh_schedule" "test" {
  dataset_id         = powerbi_pbix.test.dataset_id
  enabled            = true
  days               = ["Monday", "Wednesday", "Friday"]
  times              = ["09:00", "17:30"]
  local_time_zone_id = "Pacific Standard Time"
  notify_option      = "MailOnFailure"
}
```

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `dataset_id` - (Required, Forces new resource) The ID for the dataset that was deployed as part of the PBIX.
* `workspace_id` - (Required, Forces new resource) Workspace ID in which the dataset was deployed.
* `days` - (Required) The list of days of the week when the schedule should refresh.
* `times` - (Required) The list of times on the day the schedule should refresh. Times should be in the format HH:00 or HH:30 i.e. Hour should be two digits and minutes must either be on the full or half hour.
* `enabled` - (Optional, Default: `true`) Determines if the scheduled refresh is enabled.
* `local_time_zone_id` - (Optional, Default: `UTC`) The name of the timezone to use. See Name of Time Zone column in [Microsoft Time Zone Index Values](https://support.microsoft.com/en-gb/help/973627/microsoft-time-zone-index-values).
* `notify_option` - (Optional, Default: `NoNotification`) The notification option when a scheduled refresh fails. Should be either `MailOnFailure` or `NoNotification`.
<!-- /docgen -->
