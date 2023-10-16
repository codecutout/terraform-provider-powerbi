package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceCapacity represents a Power BI capacity
func DataSourceCapacity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCapacityRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the capacity.",
			},
			"sku": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SKU of the capacity",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region of the capacity.",
			},
		},
	}
}

func dataSourceCapacityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	name := d.Get("name").(string)
	capacityList, err := client.GetCapacities()
	if err != nil {
		return err
	}
	var capacityObjFound bool
	var capacityObj powerbiapi.GetCapacitiesResponseItem
	if len(capacityList.Value) >= 1 {
		for _, capacityObj = range capacityList.Value {
			if capacityObj.DisplayName == name {
				capacityObjFound = true
				break
			}
		}
	}
	if capacityObjFound != true {
		return fmt.Errorf("Capacity %s not found or logged-in user doesn't have capacity admin rights", name)
	}

	d.SetId(capacityObj.ID)
	d.Set("name", capacityObj.DisplayName)
	d.Set("sku", capacityObj.SKU)
	d.Set("region", capacityObj.Region)

	return nil
}
