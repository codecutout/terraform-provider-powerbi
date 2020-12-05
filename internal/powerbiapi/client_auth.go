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

type roundTripperBearerToken struct {
	innerRoundTripper http.RoundTripper
	tenant            string
	clientID          string
	clientSecret      string
	username          string
	password          string
	tokenCache        *roundTripperBearerTokenCache
}

type roundTripperBearerTokenCache struct {
	mux   sync.Mutex
	token string
}

func (rt roundTripperBearerToken) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := *req

	if rt.tokenCache.token == "" {
		err := func() error {
			rt.tokenCache.mux.Lock()
			defer rt.tokenCache.mux.Unlock()

			if rt.tokenCache.token == "" {

				// create own http client so we dont try to add token to request to get tokens
				httpClient := cleanhttp.DefaultClient()
				httpClient.Transport = roundTripperErrorOnUnsuccessful{
					innerRoundTripper: httpClient.Transport,
				}

				token, err := getAuthToken(httpClient, rt.tenant, rt.clientID, rt.clientSecret, rt.username, rt.password)
				if err != nil {
					return err
				}
				rt.tokenCache.token = token
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	newRequest.Header.Set("Authorization", "Bearer "+rt.tokenCache.token)

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
