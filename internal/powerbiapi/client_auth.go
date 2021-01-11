package powerbiapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type bearerTokenRoundTripper struct {
	innerRoundTripper http.RoundTripper
	getToken          func(*http.Client) (string, error)
	mux               sync.Mutex
	token             string
}

func newBearerTokenRoundTripper(getToken func(*http.Client) (string, error), next http.RoundTripper) http.RoundTripper {
	return &bearerTokenRoundTripper{
		innerRoundTripper: next,
		getToken:          getToken,
	}
}

func (rt *bearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := *req

	if rt.token == "" {
		err := func() error {
			rt.mux.Lock()
			defer rt.mux.Unlock()

			if rt.token == "" {

				// create own http client so we dont try to add token to request to get tokens
				httpClient := cleanhttp.DefaultClient()
				httpClient.Transport = newErrorOnUnsuccessfulRoundTripper(httpClient.Transport)

				token, err := rt.getToken(httpClient)
				if err != nil {
					return err
				}
				rt.token = token
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	newRequest.Header.Set("Authorization", "Bearer "+rt.token)

	return rt.innerRoundTripper.RoundTrip(&newRequest)
}

func getAuthTokenWithPassword(
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

func getAuthTokenWithClientCredentials(
	httpClient *http.Client,
	tenant string,
	clientID string,
	clientSecret string,
) (string, error) {

	authURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(tenant))
	resp, err := httpClient.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(url.Values{
		"grant_type":    {"client_credentials"},
		"scope":         {"https://analysis.windows.net/powerbi/api/.default"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
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
