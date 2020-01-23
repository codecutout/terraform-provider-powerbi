package api

import (
	"fmt"
	"io"
	"net/url"
	"time"
)

// PostImportInGroupRequest represents the request to create an import in a group
type PostImportInGroupRequest struct {
	GroupID            string
	Data               io.Reader `json:"-"`
	DatasetDisplayName string
	NameConflict       string
	SkipReport         bool
	Timeout            time.Duration
}

// PostImportInGroupResponse represents the response from creating an inmport in a group
type PostImportInGroupResponse struct {
	ID string
}

// GetImportInGroupRequest represents the request for getting an import in a group
type GetImportInGroupRequest struct {
	GroupID  string
	ImportID string
}

// GetImportInGroupResponse represents the response from getting an import in a group
type GetImportInGroupResponse struct {
	ID              string
	ImportState     string
	CreatedDateTime time.Time
	UpdatedDateTime time.Time
	Name            string
	ConnectionType  string
	Source          string
	Datasets        []GetImportInGroupResponseDataset
	Reports         []GetImportInGroupResponseReport
}

// GetImportInGroupResponseDataset represents the dataset from the response when getting an import in a group
type GetImportInGroupResponseDataset struct {
	ID                string
	Name              string
	WebURL            string
	TargetStorageMode string
}

// GetImportInGroupResponseReport represents the report from the response when getting an import in a group
type GetImportInGroupResponseReport struct {
	ID         string
	ReportType string
	Name       string
	WebURL     string
}

// GetImportRequest represents the request from getting an import
type GetImportRequest struct {
	ImportID string
}

// GetImportResponse represents the response from getting an import
type GetImportResponse struct {
	ID              string
	ImportState     string
	CreatedDateTime time.Time
	UpdatedDateTime time.Time
	Name            string
	ConnectionType  string
	Source          string
	Datasets        []GetImportResponseDataset
	Reports         []GetImportResponseReport
}

// GetImportResponseDataset represents the dataset from the response when getting an import
type GetImportResponseDataset struct {
	ID                string
	Name              string
	WebURL            string
	TargetStorageMode string
}

// GetImportResponseReport represents the report from the response when getting an import
type GetImportResponseReport struct {
	ID         string
	ReportType string
	Name       string
	WebURL     string
}

// PostImportInGroup creates an import wihtin the the specified group
func (client *Client) PostImportInGroup(request PostImportInGroupRequest) (*PostImportInGroupResponse, error) {

	queryParams := url.Values{}
	if request.DatasetDisplayName != "" {
		queryParams.Add("datasetDisplayName", request.DatasetDisplayName)
	}
	if request.NameConflict != "" {
		queryParams.Add("nameConflict", request.NameConflict)
	}

	var respObj PostImportInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/imports?%s", url.PathEscape(request.GroupID), queryParams.Encode())
	err := client.doMultipartJSON("POST", url, request.Data, &respObj)

	return &respObj, err
}

// WaitForImportToSucceed waits until the specified import
func (client *Client) WaitForImportToSucceed(importID string, timeout time.Duration) (*GetImportResponse, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	started := time.Now()
	for {
		im, err := client.GetImport(GetImportRequest{
			ImportID: importID,
		})
		if err != nil {
			return nil, err
		}

		if im.ImportState == "Succeeded" {
			return im, nil
		} else if im.ImportState != "Publishing" {
			return im, fmt.Errorf("Import completed with invalid state '%s'", im.ImportState)
		}

		now := <-ticker.C
		if now.Sub(started) > timeout {
			return nil, fmt.Errorf("Timed out waiting for import to complete. Import taking longer than %v seconds", timeout.Seconds())
		}
	}
}

// GetImportInGroup returns the import found within a group
func (client *Client) GetImportInGroup(request GetImportInGroupRequest) (*GetImportInGroupResponse, error) {

	var respObj GetImportInGroupResponse
	url := fmt.Sprintf(
		"https://api.powerbi.com/v1.0/myorg/groups/%s/imports/%s",
		url.PathEscape(request.GroupID),
		url.PathEscape(request.ImportID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// GetImport returns the import details
func (client *Client) GetImport(request GetImportRequest) (*GetImportResponse, error) {

	var respObj GetImportResponse
	url := fmt.Sprintf(
		"https://api.powerbi.com/v1.0/myorg/imports/%s",
		url.PathEscape(request.ImportID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}
