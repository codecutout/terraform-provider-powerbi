package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCapacity_basic(t *testing.T) {
	var capacityName = "Premium Per User - Reserved" // Using Capacity that always exists

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "powerbi_capacity" "test" {
					name = "%s"
				}
				`, capacityName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerbi_capacity.test", "name", capacityName),
					resource.TestCheckResourceAttrSet("data.powerbi_capacity.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_capacity.test", "sku"),
					resource.TestCheckResourceAttrSet("data.powerbi_capacity.test", "region"),
				),
			},
		},
	})
}
