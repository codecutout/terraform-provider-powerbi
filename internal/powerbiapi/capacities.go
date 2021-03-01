package powerbiapi

import (
	"fmt"
	"net/url"
)

// GroupAssignToCapacityRequest represents the request for Assigning capacity to group API.
type GroupAssignToCapacityRequest struct {
	CapacityID string `json:"capacityId"`
}

// GetCapacitiesResponse represents the response of get capacities response API.
type GetCapacitiesResponse struct {
	Value []GetCapacitiesResponseItem
}

//GetCapacitiesResponseItem represents the response object of each capacity of get capacities response API.
type GetCapacitiesResponseItem struct {
	ID                      string
	DisplayName             string
	Admins                  []CapacityAdmins
	SKU                     string
	State                   string
	Region                  string
	CapacityUserAccessRight string
}

//CapacityAdmins represents the list of capacity admins.
type CapacityAdmins string

// GroupAssignToCapacity assigns capcity to a workspace
func (client *Client) GroupAssignToCapacity(groupID string, request GroupAssignToCapacityRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/AssignToCapacity", url.PathEscape(groupID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

// GetCapacities Returns a list of capacities the user has access to.
func (client *Client) GetCapacities() (*GetCapacitiesResponse, error) {
	var respObj GetCapacitiesResponse
	err := client.doJSON("GET", "https://api.powerbi.com/v1.0/myorg/capacities", nil, &respObj)

	return &respObj, err
}
