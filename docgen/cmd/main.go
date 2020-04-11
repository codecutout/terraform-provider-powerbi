package main

import (
	"github.com/codecutout/terraform-provider-powerbi/docgen"
	"github.com/codecutout/terraform-provider-powerbi/powerbi"
)

func main() {
	docgen.PopulateTerraformDocs("./docs", "powerbi", powerbi.Provider())
}
