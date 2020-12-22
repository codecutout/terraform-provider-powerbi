# PBIX Resource
`powerbi_pbix` represents a PBIX upload in Power BI. 

Although from the user perspective this is a single resource, internally a PBIX upload generates the following
* Import object - identifeid with `id`
* Dataset object - identified with `dataset_id`
* Report object - identified with `report_id`

## Example Usage

### Datasource
```hcl
resource "powerbi_pbix" "mypbix" {
	workspace_id = "470b0d57-1f23-4332-a16f-9235bd174318"
	name = "My PBIX"
	source = "./my-pbix.pbix"
	source_hash = "${filemd5("./my-pbix.pbix")}"
	datasource {
		type = "OData"
		url = "https://services.odata.org/V3/(S(kbiqo1qkby04vnobw0li0fcp))/OData/OData.svc"
		original_url = "https://services.odata.org/V3/OData/OData.svc"
	}
}
```


### Parameters
```hcl
resource "powerbi_pbix" "mypbix" {
	workspace_id = "470b0d57-1f23-4332-a16f-9235bd174318"
	name = "My PBIX"
	source = "./my-pbix.pbix"
	source_hash = "${filemd5(".my-pbix.pbix")}"
	parameter {
		name = "UrlParam"
		value = "https://test-data.com/source"
	}
	parameter {
		name = "Filter"
		value = "Blue"
	}
}
```
## Argument Reference
The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required, Forces new resource) Name of the PBIX. This will be used as the name for the report and dataset.
* `workspace_id` - (Required, Forces new resource) Workspace ID in which the PBIX will be added.
* `source` - (Required) An absolute path to a PBIX file on the local system.
* `datasource` - (Optional) Datasources to be reconfigured after deploying the PBIX dataset. Changing this value will require reuploading the PBIX. Any datasource updated will not be tracked. A [`datasource`](#prop-datasource) block is defined below.
* `parameter` - (Optional) Parameters to be configured on the PBIX dataset. These can be updated wihtout requiring reuploading the PBIX. Any parameters not mentioned will not be tracked or updated. A [`parameter`](#prop-parameter) block is defined below.
* `source_hash` - (Optional) Used to trigger updates. The only meaningful value is `${filemd5("path/to/file")}`.
---
<a id="prop-datasource"></a>A `datasource` block supports the following:
* `database` - (Optional) The database name, if applicable for the type of datasource.
* `original_database` - (Optional) The database name as configured in the PBIX, if applicable for the type of datasource This will be the value replaced with the value in the 'databsase' field.
* `original_server` - (Optional) The server name as configured in the PBIX, if applicable for the type of datasource. This will be the value replaced with the value in the 'server' field.
* `original_url` - (Optional) The service URL as configured in the PBIX, if applicable for the type of datasource. This will be the value replaced with the value in the 'url' field.
* `server` - (Optional) The server name, if applicable for the type of datasource.
* `type` - (Optional) The type of datasource. For example web, sql.
* `url` - (Optional) The service URL, if applicable for the type of datasource.
---
<a id="prop-parameter"></a>A `parameter` block supports the following:
* `name` - (Required) The parameter name.
* `value` - (Required) The parameter value.
<!-- /docgen -->

## Attributes Reference
The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the import.
<!-- docgen:ComputedParameters -->
* `dataset_id` - (Optional) The ID for the dataset that was deployed as part of the PBIX.
* `report_id` - (Optional) The ID for the report that was deployed as part of the PBIX.
<!-- /docgen -->