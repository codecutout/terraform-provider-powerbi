# Power BI Provider
The PowerBI provider can be used to configure Power BI resources using the Power BI REST API
## Authentication
To use Power BI REST APIs the following need to be configured
1. A user with admin permission within Power BI. A user can be created by a domain owner, then Power BI admin permissions can be assigned via Office 365 admin center (link to the Office 365 admin center can be found within Power BI admin portal under users)
1. An Azure Active Directory App registration with delegate permissions on 'Power BI Service"' for 'Content.Create'. Easiest way to configure this is via https://dev.powerbi.com/apps 

## Example Usage
``` terraform
provider "powerbi" {
	tenant_id = "1c4cc30c-271e-47f2-891e-fef13f035bc7"
	client_id = "f9ad3042-a969-4a31-826e-856d238df3b1"
	client_secret = "u94lE93qfJSJRTEGs@Pgs]]RZzM]V?bE"
	username = "powerbiapp@mycompany.com
	password = "pass@word1!"
}
```

## Argument Reference
The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `password` - (Required) The password for the a Power BI user to use for performing Power BI REST API operations. Power BI only supports delegate permissions so a real user must be specified. This can also be sourced from the `POWERBI_PASSWORD` Environment Variable.
* `tenant_id` - (Required) The Tenant ID for the tenant which contains the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_TENANT_ID` Environment Variable.
* `client_id` - (Required) Also called Application ID. The Client ID for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_ID` Environment Variable.
* `client_secret` - (Required) Also called Application Secret. The Client Secret for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_SECRET` Environment Variable.
* `username` - (Required) The username for the a Power BI user to use for performing Power BI REST API operations. Power BI only supports delegate permissions so a real user must be specified. This can also be sourced from the `POWERBI_USERNAME` Environment Variable.
<!-- /docgen -->