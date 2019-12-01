package api

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client allows calling the Power BI service
type Client struct {
	*http.Client
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type roundTripperDefaultHeaders struct {
	rt     http.RoundTripper
	header map[string]string
}

type roundTripperErrorOnUnsuccessful struct {
	rt http.RoundTripper
}

func (h roundTripperDefaultHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.header {

		// Only set header if the request has not set it explicitly
		if existing := req.Header.Get(k); existing == "" {
			req.Header.Set(k, v)
		}
	}

	return h.rt.RoundTrip(req)
}

func (h roundTripperErrorOnUnsuccessful) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := h.rt.RoundTrip(req)

	if err != nil || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return resp, err
	}

	// Unsuccessful status code
	var reqBody []byte
	if req.ContentLength > 0 {
		resetBody, _ := req.GetBody()
		reqBody, _ = ioutil.ReadAll(resetBody)
	}

	var respBody []byte
	if resp.Body != http.NoBody {
		respBody, _ = ioutil.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("returned with status '%s'. %v", resp.Status, map[string]string{
		"request":  string(reqBody),
		"response": string(respBody),
	})

}

//NewClient creates a Power BI REST API client
func NewClient(tenant string, clientID string, clientSecret string, username string, password string) (*Client, error) {

	httpClient := cleanhttp.DefaultClient()

	authToken, err := getAuthToken(httpClient, tenant, clientID, clientSecret, username, password)
	if err != nil {
		return nil, err
	}

	// add default header for all future requests
	httpClient.Transport = roundTripperErrorOnUnsuccessful{
		roundTripperDefaultHeaders{
			httpClient.Transport,
			map[string]string{
				"Authorization": "Bearer " + authToken,
			},
		},
	}

	return &Client{
		httpClient,
	}, nil
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
