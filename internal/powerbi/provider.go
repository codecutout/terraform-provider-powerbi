package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider represents the powerbi terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_TENANT_ID", ""),
				Description: "The Tenant ID for the tenant which contains the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_TENANT_ID` Environment Variable",
			},
			"grant_type": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_GRANT_TYPE", ""),
				Description: "The grant type for the a Power BI to define token method using username and password or Application id with admin grants. This can also be sourced from the `POWERBI_GRANT_TYPE` Environment Variable",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CLIENT_ID", ""),
				Description: "Also called Application ID. The Client ID for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_ID` Environment Variable",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CLIENT_SECRET", ""),
				Description: "Also called Application Secret. The Client Secret for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_SECRET` Environment Variable",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_USERNAME", ""),
				Description: "The username for the a Power BI user to use for performing Power BI REST API operations. Power BI only supports delegate permissions so a real user must be specified. This can also be sourced from the `POWERBI_USERNAME` Environment Variable",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_PASSWORD", ""),
				Description: "The password for the a Power BI user to use for performing Power BI REST API operations. Power BI only supports delegate permissions so a real user must be specified. This can also be sourced from the `POWERBI_PASSWORD` Environment Variable",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"powerbi_workspace":        ResourceWorkspace(),
			"powerbi_pbix":             ResourcePBIX(),
			"powerbi_refresh_schedule": ResourceRefreshSchedule(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return powerbiapi.NewClient(
		d.Get("tenant_id").(string),
		d.Get("grant_type").(string),
		d.Get("client_id").(string),
		d.Get("client_secret").(string),
		d.Get("username").(string),
		d.Get("password").(string),
	)
}
