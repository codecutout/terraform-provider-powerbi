package main

import (
	"github.com/codecutout/terraform-provider-powerbi/powerbi"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return powerbi.Provider()
		},
	})
}
