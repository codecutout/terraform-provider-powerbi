package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// UpdateGroupAsAdminRequest represents the request to the UpdateGroupAsAdmin API
type UpdateGroupAsAdminRequest struct {
	GroupID string `json:"-"`
	Name    string `json:"name"`
}

// UpdateGroupAsAdmin updates a workspace
func (client *Client) UpdateGroupAsAdmin(request UpdateGroupAsAdminRequest) error {

	reqData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/admin/groups/%s", url.PathEscape(request.GroupID))
	httpRequest, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqData))
	if err != nil {
		return err
	}
	httpRequest.Header.Set("content-type", "application/json")

	_, err = client.Do(httpRequest)
	return err
}
