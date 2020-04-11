package api

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// Client allows calling the Power BI service
type Client struct {
	*http.Client
}

//NewClient creates a Power BI REST API client
func NewClient(tenant string, clientID string, clientSecret string, username string, password string) (*Client, error) {

	httpClient := cleanhttp.DefaultClient()

	// add default header for all future requests
	httpClient.Transport = roundTripperBearerToken{
		innerRoundTripper: roundTripperErrorOnUnsuccessful{httpClient.Transport},
		tenant:            tenant,
		clientID:          clientID,
		clientSecret:      clientSecret,
		username:          username,
		password:          password,
		tokenCache:        &roundTripperBearerTokenCache{},
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

func (client *Client) doMultipartJSON(method string, url string, body io.Reader, response interface{}) error {

	httpRequest, err := newMultipartRequst(method, url, body)
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

func newMultipartRequst(method string, url string, reader io.Reader) (*http.Request, error) {

	// Create multipart writer
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	partWriter, err := writer.CreatePart(textproto.MIMEHeader{})
	if err != nil {
		return nil, err
	}

	// Copy everything from reader in writer
	if _, err = io.Copy(partWriter, reader); err != nil {
		return nil, err
	}
	writer.Close()

	// Create the request from our buffer
	req, err := http.NewRequest(method, url, &buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
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
