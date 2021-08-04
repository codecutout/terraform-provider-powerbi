package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceWorkspace represents a Power BI workspace
func ResourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: createWorkspace,
		Read:   readWorkspace,
		Update: updateWorkspace,
		Delete: deleteWorkspace,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the workspace.",
			},
			"capacity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Capacity ID to be assigned to workspace.",
			},
		},
	}
}

func createWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	capacityID := d.Get("capacity_id").(string)

	resp, err := client.CreateGroup(powerbiapi.CreateGroupRequest{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	if capacityID != "" {
		err := assignToCapacity(d, meta)
		if err != nil {
			return err
		}
	}

	return readWorkspace(d, meta)
}

func readWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	workspace, err := client.GetGroup(d.Id())
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

func updateWorkspace(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("capacity_id") {
		if capacityID := d.Get("capacity_id").(string); capacityID == "" {
			d.Set("capacity_id", "00000000-0000-0000-0000-000000000000")
		}

		err := assignToCapacity(d, meta)
		if err != nil {
			return err
		}
	}

	return readWorkspace(d, meta)
}

func deleteWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	return client.DeleteGroup(d.Id())
}

func assignToCapacity(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	capacityID := d.Get("capacity_id").(string)
	if capacityID != "00000000-0000-0000-0000-000000000000" {
		var capacityObjFound bool

		capacityList, err := client.GetCapacities()
		if err != nil {
			return err
		}

		if len(capacityList.Value) >= 1 {
			for _, capacityObj := range capacityList.Value {
				if capacityObj.ID == capacityID {
					capacityObjFound = true
				}
			}
		}
		if capacityObjFound != true {
			return fmt.Errorf("Capacity id %s not found or logged-in user doesn't have capacity admin rights", capacityID)
		}
	}

	err := client.GroupAssignToCapacity(d.Id(), powerbiapi.GroupAssignToCapacityRequest{
		CapacityID: capacityID,
	})
	if err != nil {
		return err
	}

	return nil
}
