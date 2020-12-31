package powerbiapi

import (
	"fmt"
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

// GetGroupResponse represents the details when getting an individual group
type GetGroupResponse struct {
	ID                    string
	IsOnDedicatedCapacity bool
	Name                  string
}

//GetGroupUsersResponse represents list of users that have access to the specified workspace.
type GetGroupUsersResponse struct {
	Value []GetGroupUsersResponseItem
}

//GetGroupUsersResponseItem represents a single user details.
type GetGroupUsersResponseItem struct {
	DisplayName          string
	EmailAddress         string
	GroupUserAccessRight string
	Identifier           string
	PrincipalType        string
}

//AddGroupUserRequest represents details when adding a group user.
type AddGroupUserRequest struct {
	DisplayName          string `json:"displayName"`
	EmailAddress         string `json:"emailAddress"`
	GroupUserAccessRight string `json:"groupUserAccessRight"`
	Identifier           string `json:"identifier"`
	PrincipalType        string `json:"principalType"`
}

//UpdateGroupUserRequest represents details when updating a group user.
type UpdateGroupUserRequest struct {
	DisplayName          string `json:"displayName"`
	EmailAddress         string `json:"emailAddress"`
	GroupUserAccessRight string `json:"groupUserAccessRight"`
	Identifier           string `json:"identifier"`
	PrincipalType        string `json:"principalType"`
}

// CreateGroup creates new workspace
func (client *Client) CreateGroup(request CreateGroupRequest) (*CreateGroupResponse, error) {

	var respObj CreateGroupResponse
	err := client.doJSON("POST", "https://api.powerbi.com/v1.0/myorg/groups?workspaceV2=True", request, &respObj)
	return &respObj, err
}

// GetGroups returns a list of workspaces the user has access to.
func (client *Client) GetGroups(filter string, top int, skip int) (*GetGroupsResponse, error) {

	queryParams := url.Values{}
	if filter != "" {
		queryParams.Add("$filter", filter)
	}
	if top > 0 {
		queryParams.Add("$top", strconv.Itoa(top))
	}
	if skip > 0 {
		queryParams.Add("$skip", strconv.Itoa(skip))
	}

	var respObj GetGroupsResponse
	err := client.doJSON("GET", "https://api.powerbi.com/v1.0/myorg/groups?"+queryParams.Encode(), nil, &respObj)

	return &respObj, err
}

// GetGroup returns a single workspace
func (client *Client) GetGroup(groupID string) (*GetGroupResponse, error) {

	// There is no endpoint to get a single workspace, so we will search for
	// all workspaces with a specific id
	groups, err := client.GetGroups(fmt.Sprintf("id eq '%s'", groupID), -1, 0)

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

// GetGroupbyName returns a single workspace
func (client *Client) GetGroupbyName(groupName string) (*GetGroupResponse, error) {

	// There is no endpoint to get a single workspace, so we will search for
	// all workspaces with a specific name
	groups, err := client.GetGroups(fmt.Sprintf("name eq '%s'", groupName), -1, 0)

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
func (client *Client) DeleteGroup(groupID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s", url.PathEscape(groupID))
	return client.doJSON("DELETE", url, nil, nil)
}

//GetGroupUsers Returns a list of users that have access to the specified workspace.
func (client *Client) GetGroupUsers(groupID string) (*GetGroupUsersResponse, error) {

	var respObj GetGroupUsersResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/users", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

//AddGroupUser Grants the specified user permissions to the specified workspace.
func (client *Client) AddGroupUser(groupID string, request AddGroupUserRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/users", url.PathEscape(groupID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

//UpdateGroupUser Update the specified user permissions to the specified workspace.
func (client *Client) UpdateGroupUser(groupID string, request UpdateGroupUserRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/users", url.PathEscape(groupID))
	err := client.doJSON("PUT", url, &request, nil)

	return err
}

//DeleteUserInGroup Deletes the specified user permissions from the specified workspace.
func (client *Client) DeleteUserInGroup(groupID string, userInfo string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/users/%s", url.PathEscape(groupID), url.PathEscape(userInfo))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}
