package api

import (
	"fmt"
	"io"
	"net/url"
	"time"
)

// PostImportInGroupResponse represents the response from creating an inmport in a group
type PostImportInGroupResponse struct {
	ID string
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

// GetImportsInGroupResponse represents the response from imports in a group
type GetImportsInGroupResponse struct {
	Value []GetImportsInGroupResponseItem
}

// GetImportsInGroupResponseItem represents a single response item from getting an imports in a group
type GetImportsInGroupResponseItem struct {
	ID              string
	ImportState     string
	CreatedDateTime time.Time
	UpdatedDateTime time.Time
	Name            string
	ConnectionType  string
	Source          string
	Datasets        []GetImportsInGroupResponseItemDataset
	Reports         []GetImportsInGroupResponseItemReport
}

// GetImportsInGroupResponseItemDataset represents the dataset from the response when getting a imports in a group
type GetImportsInGroupResponseItemDataset struct {
	ID                string
	Name              string
	WebURL            string
	TargetStorageMode string
}

// GetImportsInGroupResponseItemReport represents the report from the response when getting a imports in a group
type GetImportsInGroupResponseItemReport struct {
	ID         string
	ReportType string
	Name       string
	WebURL     string
}

// PostImportInGroup creates an import wihtin the the specified group
func (client *Client) PostImportInGroup(groupID string, datasetDisplayName string, nameConflict string, requestData io.Reader) (*PostImportInGroupResponse, error) {

	queryParams := url.Values{}
	if datasetDisplayName != "" {
		queryParams.Add("datasetDisplayName", datasetDisplayName)
	}
	if nameConflict != "" {
		queryParams.Add("nameConflict", nameConflict)
	}

	var respObj PostImportInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/imports?%s", url.PathEscape(groupID), queryParams.Encode())
	err := client.doMultipartJSON("POST", url, requestData, &respObj)

	return &respObj, err
}

// WaitForImportToSucceed waits until the specified import
func (client *Client) WaitForImportToSucceed(importID string, timeout time.Duration) (*GetImportResponse, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	started := time.Now()
	for {
		im, err := client.GetImport(importID)
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
func (client *Client) GetImportInGroup(groupID string, importID string) (*GetImportInGroupResponse, error) {

	var respObj GetImportInGroupResponse
	url := fmt.Sprintf(
		"https://api.powerbi.com/v1.0/myorg/groups/%s/imports/%s",
		url.PathEscape(groupID),
		url.PathEscape(importID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// GetImportsInGroup returns the imports found within a group
func (client *Client) GetImportsInGroup(groupID string) (*GetImportsInGroupResponse, error) {

	var respObj GetImportsInGroupResponse
	url := fmt.Sprintf(
		"https://api.powerbi.com/v1.0/myorg/groups/%s/imports",
		url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// GetImport returns the import details
func (client *Client) GetImport(importID string) (*GetImportResponse, error) {

	var respObj GetImportResponse
	url := fmt.Sprintf(
		"https://api.powerbi.com/v1.0/myorg/imports/%s",
		url.PathEscape(importID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}
