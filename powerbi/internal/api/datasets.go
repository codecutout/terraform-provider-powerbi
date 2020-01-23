package api

import (
	"fmt"
	"net/url"
)

// GetDatasetsInGroupRequest represents the request to get datasets in a group.
type GetDatasetsInGroupRequest struct {
	GroupID string
}

// GetDatasetsInGroupResponse represents the details when getting a datasets in a group.
type GetDatasetsInGroupResponse struct {
	Value []GetDatasetsInGroupResponseItem
}

// GetDatasetsInGroupResponseItem represents a single dataset
type GetDatasetsInGroupResponseItem struct {
	ID                               string
	Name                             string
	AddRowsAPIEnabled                bool
	ConfiguredBy                     string
	IsRefreshable                    bool
	IsEffectiveIdentityRequired      bool
	IsEffectiveIdentityRolesRequired bool
	TargetStorageMode                string
}

// DeleteDatasetRequest represents the request to delete a dataset
type DeleteDatasetRequest struct {
	DatasetID string
}

// GetDatasetsInGroup returns a list of datasets within the specified group.
func (client *Client) GetDatasetsInGroup(request GetDatasetsInGroupRequest) (*GetDatasetsInGroupResponse, error) {

	var respObj GetDatasetsInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets", url.PathEscape(request.GroupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// DeleteDataset deletes a dataset.
func (client *Client) DeleteDataset(request DeleteDatasetRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s", url.PathEscape(request.DatasetID))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}
