package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
)

// CreateGroupRequest represents the request for the CreateGroup API
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// CreateGroupResponse represents the response from the CreateGroup API
type CreateGroupResponse struct {
	ID                    string
	IsOnDedicatedCapacity bool
	Name                  string
}

// GetGroupsRequest represents the request to the GetGroups API
type GetGroupsRequest struct {
	Filter string
	Top    int
	Skip   int
}

// GetGroupsResponse represents the response from the GetGroups API
type GetGroupsResponse struct {
	Value []GetGroupsResponseItem
}

// GetGroupsResponseItem represents an item returned within GetGroupsResponse
type GetGroupsResponseItem struct {
	ID                    string
	IsOnDedicatedCapacity bool
	Name                  string
}

// GetGroupRequest represents the request to get an individual group
type GetGroupRequest struct {
	GroupID string
}

// GetGroupResponse represents the details when getting an individual group
type GetGroupResponse struct {
	ID                    string
	IsOnDedicatedCapacity bool
	Name                  string
}

// DeleteGroupRequest represents the request to the DeleteGroup API
type DeleteGroupRequest struct {
	GroupID string
}

// CreateGroup creates new workspace
func (client *Client) CreateGroup(request CreateGroupRequest) (*CreateGroupResponse, error) {

	resp, err := client.DoJSONRequest("POST", "https://api.powerbi.com/v1.0/myorg/groups?workspaceV2=True", request)
	if err != nil {
		return nil, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respDataObj CreateGroupResponse
	err = json.Unmarshal(respData, &respDataObj)
	return &respDataObj, err
}

// GetGroups returns a list of workspaces the user has access to.
func (client *Client) GetGroups(request GetGroupsRequest) (*GetGroupsResponse, error) {

	queryParams := url.Values{}
	if request.Filter != "" {
		queryParams.Add("$filter", request.Filter)
	}
	if request.Top != 0 {
		queryParams.Add("$top", strconv.Itoa(request.Top))
	}
	if request.Skip != 0 {
		queryParams.Add("$skip", strconv.Itoa(request.Skip))
	}

	resp, err := client.Get("https://api.powerbi.com/v1.0/myorg/groups?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respDataObj GetGroupsResponse
	err = json.Unmarshal(respData, &respDataObj)
	return &respDataObj, err
}

// GetGroup returns a single workspace
func (client *Client) GetGroup(request GetGroupRequest) (*GetGroupResponse, error) {

	// There is no endpoint to get a single workspace, so we will search for
	// all workspaces with a specific id
	groups, err := client.GetGroups(GetGroupsRequest{
		Filter: fmt.Sprintf("id eq '%s'", request.GroupID),
	})

	if err != nil {
		return nil, err
	}

	if len(groups.Value) == 0 {
		return nil, nil
	}

	singleGroup := &groups.Value[0]
	return &GetGroupResponse{
		ID:                    singleGroup.ID,
		IsOnDedicatedCapacity: singleGroup.IsOnDedicatedCapacity,
		Name:                  singleGroup.Name,
	}, nil
}

// DeleteGroup deletes a workspace
func (client *Client) DeleteGroup(request DeleteGroupRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s", url.PathEscape(request.GroupID))
	_, err := client.DoJSONRequest("DELETE", url, nil)
	return err
}
