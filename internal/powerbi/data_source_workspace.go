package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceWorkspace represents a Power BI workspace
func DataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the workspace.",
			},
			"capacity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Capacity ID to be assigned to workspace.",
			},
		},
	}
}

func dataSourceWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	name := d.Get("name").(string)
	workspace, err := client.GetGroupByName(name)
	if err != nil {
		return err
	}

	if workspace == nil {
		d.SetId("")
	} else {
		d.SetId(workspace.ID)
		d.Set("name", workspace.Name)
		if workspace.IsOnDedicatedCapacity {
			d.Set("capacity_id", workspace.CapacityID)
		} else {
			d.Set("capacity_id", "")
		}
	}

	return nil
}
