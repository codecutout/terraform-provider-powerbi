package powerbiapi

import (
	"fmt"
	"net/url"
)

// UpdateGroupAsAdminRequest represents the request to the UpdateGroupAsAdmin API
type UpdateGroupAsAdminRequest struct {
	Name string `json:"name"`
}

// UpdateGroupAsAdmin updates a workspace
func (client *Client) UpdateGroupAsAdmin(groupID string, request UpdateGroupAsAdminRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/admin/groups/%s", url.PathEscape(groupID))
	return client.doJSON("PATCH", url, request, nil)
}
