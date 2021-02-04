package powerbiapi

import (
	"fmt"
	"net/url"
)

// PostDatasetInGroupRequest represents the request to create a push dataset
type PostDatasetInGroupRequest struct {
	Name          string                                  `json:"name,omitempty"`
	DefaultMode   string                                  `json:"defaultMode,omitempty"`
	Tables        []PostDatasetInGroupRequestTable        `json:"tables,omitempty"`
	Relationships []PostDatasetInGroupRequestRelationship `json:"relationships,omitempty"`
}

// PostDatasetInGroupRequestTable represents a table in the request to create a push dataset
type PostDatasetInGroupRequestTable struct {
	Name     string                                  `json:"name,omitempty"`
	Columns  []PostDatasetInGroupRequestTableColumn  `json:"columns,omitempty"`
	Measures []PostDatasetInGroupRequestTableMeasure `json:"measures,omitempty"`
}

// PostDatasetInGroupRequestTableColumn represents a table column in the request to create a push dataset
type PostDatasetInGroupRequestTableColumn struct {
	Name         string `json:"name,omitempty"`
	DataType     string `json:"dataType,omitempty"`
	FormatString string `json:"formatString,omitempty"`
}

// PostDatasetInGroupRequestTableMeasure represents a table measure in the request to create a push dataset
type PostDatasetInGroupRequestTableMeasure struct {
	Name       string `json:"name,omitempty"`
	Expression string `json:"expression,omitempty"`
}

// PostDatasetInGroupRequestRelationship represents a relationship in the request to create a push dataset
type PostDatasetInGroupRequestRelationship struct {
	Name                   string `json:"name,omitempty"`
	FromColumn             string `json:"fromColumn,omitempty"`
	FromTable              string `json:"fromTable,omitempty"`
	ToColumn               string `json:"toColumn,omitempty"`
	ToTable                string `json:"toTable,omitempty"`
	CrossFilteringBehavior string `json:"crossFilteringBehavior,omitempty"`
}

// PostDatasetInGroupResponse represents the details when post a dataset in a group.
type PostDatasetInGroupResponse struct {
	ID   string
	Name string
}

// GetTablesResponse represents the response to a request to get tables
type GetTablesResponse struct {
	Value []GetTablesResponseTable
}

// GetTablesResponseTable represents a table from the response to a request to get tables
type GetTablesResponseTable struct {
	Name string
}

// PutTableInGroupRequest represents the request to update a table
type PutTableInGroupRequest struct {
	Name     string                               `json:"name,omitempty"`
	Columns  []PutTableInGroupRequestTableColumn  `json:"columns,omitempty"`
	Measures []PutTableInGroupRequestTableMeasure `json:"measures,omitempty"`
}

// PutTableInGroupRequestTableColumn represents a table column in the request to update a table
type PutTableInGroupRequestTableColumn struct {
	Name         string `json:"name,omitempty"`
	DataType     string `json:"dataType,omitempty"`
	FormatString string `json:"formatString,omitempty"`
}

// PutTableInGroupRequestTableMeasure represents a table measure in the request to update a table
type PutTableInGroupRequestTableMeasure struct {
	Name       string `json:"name,omitempty"`
	Expression string `json:"expression,omitempty"`
}

// PostRowsInGroupRequest represents the request to post rows into a push dataset
type PostRowsInGroupRequest struct {
	Rows []map[string]interface{}
}

// PostDatasetInGroup creates a dataset within the specified group.
func (client *Client) PostDatasetInGroup(groupID string, defaultRetentionPolicy string, request PostDatasetInGroupRequest) (*PostDatasetInGroupResponse, error) {

	queryParams := url.Values{}
	if defaultRetentionPolicy != "" {
		queryParams.Add("defaultRetentionPolicy", defaultRetentionPolicy)
	}

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets?%s",
		url.PathEscape(groupID),
		queryParams.Encode())

	var respObj PostDatasetInGroupResponse
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}

// GetTables gets the tables in a push dataset.
func (client *Client) GetTables(datasetID string) (*GetTablesResponse, error) {

	var respObj GetTablesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/tables", url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// PutTableInGroup updates the metadata and schema for the specified table, within the specified dataset, from the specified workspace.
func (client *Client) PutTableInGroup(groupID string, datasetID string, tableName string, request PutTableInGroupRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/tables/%s",
		url.PathEscape(groupID),
		url.PathEscape(datasetID),
		url.PathEscape(tableName))

	return client.doJSON("PUT", url, &request, nil)
}

// PostRowsInGroup posts rows into a table in a dataset in a group.
func (client *Client) PostRowsInGroup(groupID string, datasetID string, tableName string, request PostRowsInGroupRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/tables/%s/rows",
		url.PathEscape(groupID),
		url.PathEscape(datasetID),
		url.PathEscape(tableName))
	return client.doJSON("POST", url, &request, nil)
}
