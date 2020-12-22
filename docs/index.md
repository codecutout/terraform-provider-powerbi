# Power BI Provider
The PowerBI provider can be used to configure Power BI resources using the [Power BI REST API](https://docs.microsoft.com/en-us/rest/api/power-bi/)

-> See the [authentication guide](guides/authentication.md) for details on how to generate authetnication details

## Example Usage
```hcl
provider "powerbi" {
	tenant_id = "1c4cc30c-271e-47f2-891e-fef13f035bc7"
	client_id = "f9ad3042-a969-4a31-826e-856d238df3b1"
	client_secret = "u94lE93qfJSJRTEGs@Pgs]]RZzM]V?bE"
	username = "powerbiapp@mycompany.com"
	password = "pass@word1!"
}
```

## Argument Reference
The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `client_id` - (Required) Also called Application ID. The Client ID for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_ID` Environment Variable.
* `client_secret` - (Required) Also called Application Secret. The Client Secret for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_SECRET` Environment Variable.
* `tenant_id` - (Required) The Tenant ID for the tenant which contains the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_TENANT_ID` Environment Variable.
* `password` - (Optional) The password for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_PASSWORD` Environment Variable.
* `username` - (Optional) The username for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_USERNAME` Environment Variable.
<!-- /docgen -->