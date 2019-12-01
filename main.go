package main

import (
	"github.com/alex-davies/terraform-provider-powerbi/powerbi"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: powerbi.Provider})
}
