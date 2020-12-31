package powerbiapi

//RefreshUserPermissions Refreshes user permissions in Power BI.
func (client *Client) RefreshUserPermissions() error {
	err := client.doJSON("POST", "https://api.powerbi.com/v1.0/myorg/RefreshUserPermissions", nil, nil)

	return err
}
