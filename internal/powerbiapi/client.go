package powerbiapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// Client allows calling the Power BI service
type Client struct {
	*http.Client
}

//NewClientWithPasswordAuth creates a Power BI REST API client using password authentication with delegated permissions
func NewClientWithPasswordAuth(tenant string, clientID string, clientSecret string, username string, password string) (*Client, error) {
	return newClient(func(httpClient *http.Client) (string, error) {
		return getAuthTokenWithPassword(httpClient, tenant, clientID, clientSecret, username, password)
	})
}

//NewClientWithClientCredentialAuth creates a Power BI REST API client using client credentials with application permissions
func NewClientWithClientCredentialAuth(tenant string, clientID string, clientSecret string) (*Client, error) {

	return newClient(func(httpClient *http.Client) (string, error) {
		return getAuthTokenWithClientCredentials(httpClient, tenant, clientID, clientSecret)
	})
}

func newClient(getAuthToken func(httpClient *http.Client) (string, error)) (*Client, error) {

	// PowerBI has lots of intermittant TLS handshake issues, these settings
	// seem to reduce the amount of issues encountered
	defaultTransport := cleanhttp.DefaultPooledTransport()
	defaultTransport.TLSHandshakeTimeout = 60 * time.Second
	defaultTransport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// auth
	httpClient := &http.Client{
		Transport: newBearerTokenRoundTripper(
			getAuthToken,
			// error
			newErrorOnUnsuccessfulRoundTripper(
				// retry internal server error
				// this is crazy that we have to do this, but the API intermittently returns internal server errors
				newRetryInternalServerErrorRoundTripper(
					// retry too many requests
					newRetryTooManyRequestsRoundTripper(
						// actual call
						defaultTransport,
					),
				),
			),
		),
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

	httpRequest, err := newMultipartRequest(method, url, body)
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

func newMultipartRequest(method string, url string, reader io.Reader) (*http.Request, error) {

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
