package powerbi

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceGroupUsers represents user management in Power BI workspace.
func ResourceGroupUsers() *schema.Resource {
	return &schema.Resource{
		Create: addGroupUser,
		Read:   readGroupUser,
		Update: updateGroupUser,
		Delete: deleteGroupUser,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "Workspace ID in which the PBIX will be added.",
				Required:    true,
			},
			"group_user_access_right": {
				Type:        schema.TypeString,
				Description: "User access level to workspace. Any value from Admin, Contributor, Member, None or Viewer.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					stringVal := val.(string)
					reg := regexp.MustCompile("^(Admin|Contributor|Member|None|Viewer)$")
					if !reg.MatchString(stringVal) {
						errs = append(errs, fmt.Errorf("Expected argument 'group_user_access_right' to have value one of Admin, Contributor, Member, None or Viewer. Found '%v'", stringVal))
					}
					return warns, errs
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of the principal.",
				Optional:    true,
				Computed:    true,
			},
			"email_address": {
				Type:        schema.TypeString,
				Description: "Email address of the user.",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					stringVal := val.(string)
					reg := regexp.MustCompile(".*@.*")
					if !reg.MatchString(stringVal) {
						errs = append(errs, fmt.Errorf("Expected argument 'email_address' to be like user@mailserver.com. Found '%v'", stringVal))
					}
					return warns, errs
				},
			},
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of the principal.",
				Optional:    true,
				Computed:    true,
			},
			"principal_type": {
				Type:        schema.TypeString,
				Description: "The principal type. Any value from App, Group or User.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					stringVal := val.(string)
					reg := regexp.MustCompile("^(User|App|Group)$")
					if !reg.MatchString(stringVal) {
						errs = append(errs, fmt.Errorf("Expected argument 'principal_type' to have value one of User, Group or App. Found '%v'", stringVal))
					}
					return warns, errs
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func addGroupUser(d *schema.ResourceData, meta interface{}) error {

	groupID := d.Get("workspace_id").(string)

	Identifier := d.Get("identifier").(string)
	if Identifier == "" {
		Identifier = d.Get("email_address").(string)
	}

	client := meta.(*powerbiapi.Client)
	err := client.AddGroupUser(groupID, powerbiapi.GroupUserDetails{
		GroupUserAccessRight: d.Get("group_user_access_right").(string),
		DisplayName:          d.Get("display_name").(string),
		PrincipalType:        d.Get("principal_type").(string),
		EmailAddress:         d.Get("email_address").(string),
		Identifier:           d.Get("identifier").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(groupID + "/" + Identifier)
	return nil
}

func readGroupUser(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*powerbiapi.Client)
	groupID := d.Get("workspace_id").(string)
	if groupID == "" {
		groupID = strings.SplitN(d.Id(), "/", 2)[0]
		fmt.Printf("readGroupUser groupID %s fetched from state id %s", groupID, d.Id())
	}

	Identifier := d.Get("identifier").(string)
	if Identifier == "" {
		Identifier = d.Get("email_address").(string)
	}
	if Identifier == "" {
		Identifier = strings.SplitN(d.Id(), "/", 2)[1]
		fmt.Printf("deleteGroupUser Identifier %s fetched from state id %s", Identifier, d.Id())
	}
	if Identifier == "" {
		return fmt.Errorf("Could not find user identifier")
	}

	groupUsers, err := client.GetGroupUsers(groupID)
	if err != nil {
		return err
	}

	if len(groupUsers.Value) >= 1 {
		for _, apiOUTuserObj := range groupUsers.Value {
			if apiOUTuserObj.Identifier == Identifier {
				d.Set("identifier", apiOUTuserObj.Identifier)
				d.Set("group_user_access_right", apiOUTuserObj.GroupUserAccessRight)
				d.Set("display_name", apiOUTuserObj.DisplayName)
				d.Set("email_address", apiOUTuserObj.EmailAddress)
				d.Set("principal_type", apiOUTuserObj.PrincipalType)
				d.Set("workspace_id", groupID)
			}
		}
	}

	return nil
}

func updateGroupUser(d *schema.ResourceData, meta interface{}) error {

	groupID := d.Get("workspace_id").(string)
	if groupID == "" {
		groupID = strings.SplitN(d.Id(), "/", 2)[0]
		fmt.Printf("updateGroupUser groupID %s fetched from state id %s", groupID, d.Id())
	}

	if d.HasChange("group_user_access_right") {

		client := meta.(*powerbiapi.Client)
		err := client.UpdateGroupUser(groupID, powerbiapi.GroupUserDetails{
			GroupUserAccessRight: d.Get("group_user_access_right").(string),
			DisplayName:          d.Get("display_name").(string),
			PrincipalType:        d.Get("principal_type").(string),
			EmailAddress:         d.Get("email_address").(string),
			Identifier:           d.Get("identifier").(string),
		})
		if err != nil {
			return err
		}

	}

	return readGroupUser(d, meta)

}

func deleteGroupUser(d *schema.ResourceData, meta interface{}) error {

	groupID := d.Get("workspace_id").(string)
	if groupID == "" {
		groupID = strings.SplitN(d.Id(), "/", 2)[0]
		fmt.Printf("deleteGroupUser groupID %s fetched from state id %s", groupID, d.Id())
	}

	Identity := d.Get("email_address")
	if Identity == nil {
		Identity = d.Get("identity")
	}
	if Identity == nil {
		Identity = strings.SplitN(d.Id(), "/", 2)[1]
		fmt.Printf("deleteGroupUser Identity %s fetched from state id %s", Identity, d.Id())
	}

	client := meta.(*powerbiapi.Client)
	return client.DeleteUserInGroup(groupID, Identity.(string))
}
