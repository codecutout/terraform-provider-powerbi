package api

import (
	"fmt"
	"net/url"
)

// GetReportsInGroupRequest represents the request to get reports in a group.
type GetReportsInGroupRequest struct {
	GroupID string
}

// GetReportsInGroupResponse represents the details when getting a report in a group.
type GetReportsInGroupResponse struct {
	Value []GetReportsInGroupResponseItem
}

// GetReportsInGroupResponseItem represents a single dataset
type GetReportsInGroupResponseItem struct {
	ID        string
	Name      string
	DatasetID string
	WebURL    string
	EmbedURL  string
}

// GetReportsInGroup returns a list of reports within the specified group.
func (client *Client) GetReportsInGroup(request GetReportsInGroupRequest) (*GetReportsInGroupResponse, error) {

	var respObj GetReportsInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports", url.PathEscape(request.GroupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}
