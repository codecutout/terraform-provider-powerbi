package api

import (
	"fmt"
	"net/url"
)

// UpdateGroupAsAdminRequest represents the request to the UpdateGroupAsAdmin API
type UpdateGroupAsAdminRequest struct {
	GroupID string `json:"-"`
	Name    string `json:"name"`
}

// UpdateGroupAsAdmin updates a workspace
func (client *Client) UpdateGroupAsAdmin(request UpdateGroupAsAdminRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/admin/groups/%s", url.PathEscape(request.GroupID))
	return client.doJSON("PATCH", url, request, nil)
}
