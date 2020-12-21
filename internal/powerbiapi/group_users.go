package powerbiapi

import (
	"fmt"
	"net/url"
)

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

//GroupUserDetails represents a single user details.
type GroupUserDetails struct {
	DisplayName          string `json:"displayName"`
	EmailAddress         string `json:"emailAddress"`
	GroupUserAccessRight string `json:"groupUserAccessRight"`
	Identifier           string `json:"identifier"`
	PrincipalType        string `json:"principalType"`
}

//GetGroupUsers Returns a list of users that have access to the specified workspace.
func (client *Client) GetGroupUsers(groupID string) (*GetGroupUsersResponse, error) {

	var respObj GetGroupUsersResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/users", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

//AddGroupUser Grants the specified user permissions to the specified workspace.
func (client *Client) AddGroupUser(groupID string, request GroupUserDetails) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/users", url.PathEscape(groupID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

//UpdateGroupUser Update the specified user permissions to the specified workspace.
func (client *Client) UpdateGroupUser(groupID string, request GroupUserDetails) error {
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

//RefreshUserPermissions Refreshes user permissions in Power BI.
func (client *Client) RefreshUserPermissions() error {
	err := client.doJSON("POST", "https://api.powerbi.com/v1.0/myorg/RefreshUserPermissions", nil, nil)

	return err
}
