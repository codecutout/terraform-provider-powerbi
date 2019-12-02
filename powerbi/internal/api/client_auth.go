package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type roundTripperBearerToken struct {
	innerRoundTripper http.RoundTripper
	token             string
}

func (rt roundTripperBearerToken) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := *req
	newRequest.Header.Set("Authorization", "Bearer "+rt.token)

	return rt.innerRoundTripper.RoundTrip(&newRequest)
}

func getAuthToken(
	httpClient *http.Client,
	tenant string,
	clientID string,
	clientSecret string,
	username string,
	password string,
) (string, error) {

	authURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(tenant))
	resp, err := httpClient.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(url.Values{
		"grant_type":    {"password"},
		"scope":         {"https://analysis.windows.net/powerbi/api/.default"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password},
	}.Encode()))

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("status: %d, body: %s", resp.StatusCode, data)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var dataObj tokenResponse
	err = json.Unmarshal(data, &dataObj)
	return dataObj.AccessToken, err
}
