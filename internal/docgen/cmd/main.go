package main

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/docgen"
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbi"
)

func main() {
	docgen.PopulateTerraformDocs("./docs", "powerbi", powerbi.Provider())
}
