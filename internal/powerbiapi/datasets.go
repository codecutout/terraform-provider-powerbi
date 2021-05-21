package powerbiapi

import (
	"fmt"
	"net/url"
)

// GetDatasetInGroupResponse represents the details when getting a datasets in a group.
type GetDatasetInGroupResponse struct {
	ID                               string `json:"id"`
	Name                             string
	AddRowsAPIEnabled                bool
	ConfiguredBy                     string
	IsRefreshable                    bool
	IsEffectiveIdentityRequired      bool
	IsEffectiveIdentityRolesRequired bool
	TargetStorageMode                string
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

// GetParametersInGroupResponse represents the response from get parameters
type GetParametersInGroupResponse struct {
	Value []GetParametersInGroupResponseItem
}

// GetParametersInGroupResponseItem represents a single parameter
type GetParametersInGroupResponseItem struct {
	Name         string
	Type         string
	IsRequired   bool
	CurrentValue string
}

// UpdateParametersInGroupRequest represents the request to update parameters
type UpdateParametersInGroupRequest struct {
	UpdateDetails []UpdateParametersInGroupRequestItem
}

// UpdateParametersInGroupRequestItem represents a single parameter update
type UpdateParametersInGroupRequestItem struct {
	Name     string
	NewValue string
}

// GetDatasourcesInGroupResponse represents the response from get datasources
type GetDatasourcesInGroupResponse struct {
	Value []GetDatasourcesInGroupResponseItem
}

// GetDatasourcesInGroupResponseItem represents a single datasource
type GetDatasourcesInGroupResponseItem struct {
	DatasourceID      string
	DatasourceType    string
	GatewayID         string
	Name              string
	CopnnectionString string
	ConnectionDetails GetDatasourcesInGroupResponseItemConnectionDetails
}

// GetDatasourcesInGroupResponseItemConnectionDetails represents connection details for a single datasource
type GetDatasourcesInGroupResponseItemConnectionDetails struct {
	Database *string
	Server   *string
	URL      *string
}

// UpdateDatasourcesInGroupRequest represents the request to update datasources
type UpdateDatasourcesInGroupRequest struct {
	UpdateDetails []UpdateDatasourcesInGroupRequestItem
}

// UpdateDatasourcesInGroupRequestItem represents a single datasource update
type UpdateDatasourcesInGroupRequestItem struct {
	DatasourceSelector UpdateDatasourcesInGroupRequestItemDatasourceSelector
	ConnectionDetails  UpdateDatasourcesInGroupRequestItemConnectionDetails
}

// UpdateDatasourcesInGroupRequestItemDatasourceSelector represents a query to select a datasource
type UpdateDatasourcesInGroupRequestItemDatasourceSelector struct {
	DatasourceType    string
	ConnectionDetails UpdateDatasourcesInGroupRequestItemConnectionDetails
}

// UpdateDatasourcesInGroupRequestItemConnectionDetails represents connection details for a single datasource
type UpdateDatasourcesInGroupRequestItemConnectionDetails struct {
	Database *string
	Server   *string
	URL      *string
}

// GetRefreshScheduleInGroupResponse represents the response to getting a refresh schedule
type GetRefreshScheduleInGroupResponse struct {
	Enabled         bool
	Days            []string
	Times           []string
	LocalTimeZoneID string
	NotifyOption    string
}

// UpdateRefreshScheduleInGroupRequest represents the request to update refresh schedules
type UpdateRefreshScheduleInGroupRequest struct {
	Value UpdateRefreshScheduleInGroupRequestValue `json:"value"`
}

// UpdateRefreshScheduleInGroupRequestValue represents the value section in the request tot update refresh schedules
type UpdateRefreshScheduleInGroupRequestValue struct {
	Enabled         *bool     `json:"enabled,omitempty"`
	Days            *[]string `json:"days,omitempty"`
	Times           *[]string `json:"times,omitempty"`
	LocalTimeZoneID *string   `json:"localTimeZoneId,omitempty"`
	NotifyOption    *string   `json:"notifyOption,omitempty"`
}

// GetDatasetInGroup returns a dataset within the specified group.
func (client *Client) GetDatasetInGroup(groupID string, datasetID string) (*GetDatasetInGroupResponse, error) {

	var respObj GetDatasetInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// GetDatasetsInGroup returns a list of datasets within the specified group.
func (client *Client) GetDatasetsInGroup(groupID string) (*GetDatasetsInGroupResponse, error) {

	var respObj GetDatasetsInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// DeleteDatasetInGroup deletes a dataset that exists within a group.
func (client *Client) DeleteDatasetInGroup(groupID string, datasetID string) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}

// GetParametersInGroup gets parameters in a dataset that exists within a group.
func (client *Client) GetParametersInGroup(groupID string, datasetID string) (*GetParametersInGroupResponse, error) {

	var respObj GetParametersInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/parameters", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// UpdateParametersInGroup updates parameters in a dataset that exists within a group.
func (client *Client) UpdateParametersInGroup(groupID string, datasetID string, request UpdateParametersInGroupRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/Default.UpdateParameters", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

// GetDatasourcesInGroup gets datasources in a dataset that exists within a group.
func (client *Client) GetDatasourcesInGroup(groupID string, datasetID string) (*GetDatasourcesInGroupResponse, error) {

	var respObj GetDatasourcesInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/datasources", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// UpdateDatasourcesInGroup updates datasources in a dataset that exists within a group.
func (client *Client) UpdateDatasourcesInGroup(groupID string, datasetID string, request UpdateDatasourcesInGroupRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/Default.UpdateDatasources", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

// GetRefreshScheduleInGroup gets a datasource's refresh schedule.
func (client *Client) GetRefreshScheduleInGroup(groupID string, datasetID string) (*GetRefreshScheduleInGroupResponse, error) {

	var respObj GetRefreshScheduleInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/refreshSchedule", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// UpdateRefreshScheduleInGroup updates a datasource's refresh schedule.
func (client *Client) UpdateRefreshScheduleInGroup(groupID string, datasetID string, request UpdateRefreshScheduleInGroupRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/refreshSchedule", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("PATCH", url, &request, nil)

	return err
}

// TakeOverInGroup Transfers ownership over the specified dataset to the current authorized user.
func (client *Client) TakeOverInGroup(groupID string, datasetID string) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/Default.TakeOver", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("POST", url, nil, nil)

	return err
}
