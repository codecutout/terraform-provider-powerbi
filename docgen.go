package main

import (
	"github.com/alex-davies/terraform-provider-powerbi/docgen"
	"github.com/alex-davies/terraform-provider-powerbi/powerbi"
)

func main() {
	docgen.PopulateTerraformDocs("./docs", "powerbi", powerbi.Provider())
}
