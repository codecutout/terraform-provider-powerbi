package powerbiapi

import (
	"fmt"
	"net/url"
)

// RebindReportInGroup represents the request for the RebindReportInGroup API
type RebindReportInGroupRequest struct {
	DatasetID string `json:"datasetId"`
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
func (client *Client) GetReportsInGroup(groupID string) (*GetReportsInGroupResponse, error) {

	var respObj GetReportsInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// DeleteReportInGroup deletes a report that exists within a group.
func (client *Client) DeleteReportInGroup(groupID string, reportID string) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports/%s", url.PathEscape(groupID), url.PathEscape(reportID))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}

// RebindReportInGroup rebinds the specified report from the specified group to the requested dataset.
func (client *Client) RebindReportInGroup(groupID string, reportID string, request RebindReportInGroupRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports/%s/Rebind", url.PathEscape(groupID), url.PathEscape(reportID))
	err := client.doJSON("POST", url, request, nil)

	return err
}
