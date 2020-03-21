package main

import (
	"github.com/alex-davies/terraform-provider-powerbi/powerbi"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return powerbi.Provider()
		},
	})
}
