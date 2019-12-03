package api

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/go-cleanhttp"
	"io/ioutil"
	"net/http"
)

// Client allows calling the Power BI service
type Client struct {
	*http.Client
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
		roundTripperBearerToken{
			innerRoundTripper: httpClient.Transport,
			token:             authToken,
		},
	}

	return &Client{
		httpClient,
	}, nil
}

func (client *Client) doJSON(method string, url string, body interface{}, response interface{}) error {

	httpRequest, err := newJSONRequest(method, url, body)
	if err != nil {
		return err
	}

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return err
	}

	return newJSONResponse(httpResponse, response)
}

func newJSONRequest(method string, url string, body interface{}) (*http.Request, error) {

	// if we have no body so can create a simple request
	if body == nil {
		return http.NewRequest(method, url, nil)
	}

	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("content-type", "application/json")

	return httpRequest, nil
}

func newJSONResponse(httpResponse *http.Response, response interface{}) error {
	if response == nil {
		return nil
	}

	httpResponseData, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(httpResponseData, response)
}
