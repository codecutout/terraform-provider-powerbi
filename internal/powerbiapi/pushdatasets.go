package powerbiapi

import (
	"fmt"
	"net/url"
)

// PostDatasetRequest represents the request to create a push dataset
type PostDatasetRequest struct {
	Name          string
	DefaultMode   string
	Table         []PostDatasetRequestTable
	Relationships []PostDatasetRequestRelationship
}

// PostDatasetRequestTable represents a table in the request to create a push dataset
type PostDatasetRequestTable struct {
	Name     string
	Columns  []PostDatasetRequestTableColumn
	Measures []PostDatasetRequestTableMeasure
}

// PostDatasetRequestTableColumn represents a table column in the request to create a push dataset
type PostDatasetRequestTableColumn struct {
	Name     string
	DataType string
}

// PostDatasetRequestTableMeasure represents a table measure in the request to create a push dataset
type PostDatasetRequestTableMeasure struct {
	Name       string
	Expression string
}

// PostDatasetRequestRelationship represents a relationship in the request to create a push dataset
type PostDatasetRequestRelationship struct {
	Name                   string
	FromColumn             string
	FromTable              string
	ToColumn               string
	ToTable                string
	CrossFilteringBehavior string
}

// GetTablesResponse represents the response to a request to get tables
type GetTablesResponse struct {
	Value []GetTablesResponseTable
}

// GetTablesResponseTable represents a table from the response to a request to get tables
type GetTablesResponseTable struct {
	Name string
}

// PostDatasetInGroup creates a dataset within the specified group.
func (client *Client) PostDatasetInGroup(groupID string, defaultRetentionPolicy string, request PostDatasetRequest) error {

	queryParams := url.Values{}
	if defaultRetentionPolicy != "" {
		queryParams.Add("defaultRetentionPolicy", defaultRetentionPolicy)
	}

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets?%s",
		url.PathEscape(groupID),
		queryParams.Encode())
	return client.doJSON("POST", url, &request, nil)
}

// GetTables gets the tables in a push dataset.
func (client *Client) GetTables(datasetID string) (*GetTablesResponse, error) {

	var respObj GetTablesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/tables", url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}
