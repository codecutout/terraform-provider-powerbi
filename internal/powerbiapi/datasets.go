package powerbiapi

import (
	"fmt"
	"net/url"
)

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

// GetParametersResponse represents the response from get parameters
type GetParametersResponse struct {
	Value []GetParametersResponseItem
}

// GetParametersResponseItem represents a single parameter
type GetParametersResponseItem struct {
	Name         string
	Type         string
	IsRequired   bool
	CurrentValue string
}

// UpdateParametersRequest represents the request to update parameters
type UpdateParametersRequest struct {
	UpdateDetails []UpdateParametersRequestItem
}

// UpdateParametersRequestItem represents a single parameter update
type UpdateParametersRequestItem struct {
	Name     string
	NewValue string
}

// GetDatasourcesResponse represents the response from get datasources
type GetDatasourcesResponse struct {
	Value []GetDatasourcesResponseItem
}

// GetDatasourcesResponseItem represents a single datasource
type GetDatasourcesResponseItem struct {
	DatasourceID      string
	DatasourceType    string
	GatewayID         string
	Name              string
	CopnnectionString string
	ConnectionDetails GetDatasourcesResponseItemConnectionDetails
}

// GetDatasourcesResponseItemConnectionDetails represents connection details for a single datasource
type GetDatasourcesResponseItemConnectionDetails struct {
	Database *string
	Server   *string
	URL      *string
}

// UpdateDatasourcesRequest represents the request to update datasources
type UpdateDatasourcesRequest struct {
	UpdateDetails []UpdateDatasourcesRequestItem
}

// UpdateDatasourcesRequestItem represents a single datasource update
type UpdateDatasourcesRequestItem struct {
	DatasourceSelector UpdateDatasourcesRequestItemDatasourceSelector
	ConnectionDetails  UpdateDatasourcesRequestItemConnectionDetails
}

// UpdateDatasourcesRequestItemDatasourceSelector represents a query to select a datasource
type UpdateDatasourcesRequestItemDatasourceSelector struct {
	DatasourceType    string
	ConnectionDetails UpdateDatasourcesRequestItemConnectionDetails
}

// UpdateDatasourcesRequestItemConnectionDetails represents connection details for a single datasource
type UpdateDatasourcesRequestItemConnectionDetails struct {
	Database *string
	Server   *string
	URL      *string
}

// GetRefreshScheduleResponse represents the response to getting a refresh schedule
type GetRefreshScheduleResponse struct {
	Enabled         bool
	Days            []string
	Times           []string
	LocalTimeZoneID string
	NotifyOption    string
}

// UpdateRefreshScheduleRequest represents the request to update refresh schedules
type UpdateRefreshScheduleRequest struct {
	Value UpdateRefreshScheduleRequestValue `json:"value"`
}

// UpdateRefreshScheduleRequestValue represents the value section in the request tot update refresh schedules
type UpdateRefreshScheduleRequestValue struct {
	Enabled         *bool     `json:"enabled,omitempty"`
	Days            *[]string `json:"days,omitempty"`
	Times           *[]string `json:"times,omitempty"`
	LocalTimeZoneID *string   `json:"localTimeZoneId,omitempty"`
	NotifyOption    *string   `json:"notifyOption,omitempty"`
}

// GetDatasetsInGroup returns a list of datasets within the specified group.
func (client *Client) GetDatasetsInGroup(groupID string) (*GetDatasetsInGroupResponse, error) {

	var respObj GetDatasetsInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// DeleteDataset deletes a dataset.
func (client *Client) DeleteDataset(datasetID string) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s", url.PathEscape(datasetID))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}

// GetParameters gets parameters in a dataset.
func (client *Client) GetParameters(datasetID string) (*GetParametersResponse, error) {

	var respObj GetParametersResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/parameters", url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// UpdateParameters updates parameters in a dataset.
func (client *Client) UpdateParameters(datasetID string, request UpdateParametersRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/Default.UpdateParameters", url.PathEscape(datasetID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

// GetDatasources gets datasources in a dataset.
func (client *Client) GetDatasources(datasetID string) (*GetDatasourcesResponse, error) {

	var respObj GetDatasourcesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/datasources", url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// UpdateDatasources updates datasources in a dataset.
func (client *Client) UpdateDatasources(datasetID string, request UpdateDatasourcesRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/Default.UpdateDatasources", url.PathEscape(datasetID))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

// GetRefreshSchedule gets a datasource's refresh schedule.
func (client *Client) GetRefreshSchedule(datasetID string) (*GetRefreshScheduleResponse, error) {

	var respObj GetRefreshScheduleResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/refreshSchedule", url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// UpdateRefreshSchedule updates a datasource's refresh schedule.
func (client *Client) UpdateRefreshSchedule(datasetID string, request UpdateRefreshScheduleRequest) error {

	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/datasets/%s/refreshSchedule", url.PathEscape(datasetID))
	err := client.doJSON("PATCH", url, &request, nil)

	return err
}
