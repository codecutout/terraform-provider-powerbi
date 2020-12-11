package powerbiapi

import (
	"fmt"
	"net/url"
)

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
func (client *Client) GetReportsInGroup(groupID string) (*GetReportsInGroupResponse, error) {

	var respObj GetReportsInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// DeleteReportInGroup deletes a dataset.
func (client *Client) DeleteReportInGroup(groupID string, reportID string) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports/%s", url.PathEscape(groupID), url.PathEscape(reportID))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}
