package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description: "Name of the workspace.",
			},
			"capacity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Capacity ID, which will be assigned to workspace.",
			},
		},
	}
}

func createWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	resp, err := client.CreateGroup(powerbiapi.CreateGroupRequest{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	if d.Get("capacity_id").(string) != "" {
		err := client.GroupAssignToCapacity(d.Id(), powerbiapi.GroupAssignToCapacityRequest{
			CapacityID: d.Get("capacity_id").(string),
		})
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
		d.Set("capacity_id", workspace.CapacityID)
	}

	return nil
}

func updateWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	existingWorkspace, err := client.GetGroup(d.Id())
	if err != nil {
		return err
	}

	if existingWorkspace.Name != d.Get("name").(string) {
		err := client.UpdateGroupAsAdmin(d.Id(), powerbiapi.UpdateGroupAsAdminRequest{
			Name: d.Get("name").(string),
		})
		if err != nil {
			return err
		}
	} else {
		err := client.GroupAssignToCapacity(d.Id(), powerbiapi.GroupAssignToCapacityRequest{
			CapacityID: d.Get("capacity_id").(string),
		})
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
