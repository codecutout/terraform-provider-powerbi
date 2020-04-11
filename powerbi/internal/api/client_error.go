package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPUnsuccessfulError represents an error thrown when a non 2xx response is received
type HTTPUnsuccessfulError struct {
	Request      *http.Request
	Response     *http.Response
	ErrorBody    *ErrorBody
	ErrorBodyRaw []byte
}

// ErrorResponse represents the response when the Power BI API returns errors
type ErrorResponse struct {
	Error ErrorBody
}

// ErrorBody represents the error returend in the body of Power BI API requests
type ErrorBody struct {
	Code    string
	Message string
}

type roundTripperErrorOnUnsuccessful struct {
	innerRoundTripper http.RoundTripper
}

func (err HTTPUnsuccessfulError) Error() string {

	message := fmt.Sprintf("status code '%s'", err.Response.Status)
	if err.ErrorBody != nil && err.ErrorBody.Code != "" && err.ErrorBody.Message != "" {
		message += fmt.Sprintf(" with code '%s' and message '%s'", err.ErrorBody.Code, err.ErrorBody.Message)
	} else if err.ErrorBody != nil && err.ErrorBody.Code != "" {
		message += fmt.Sprintf(" with code '%s'", err.ErrorBody.Code)
	} else if len(err.ErrorBodyRaw) > 0 {
		message += fmt.Sprintf(" with body %s", string(err.ErrorBodyRaw))
	}
	return message
}

func (h roundTripperErrorOnUnsuccessful) RoundTrip(req *http.Request) (*http.Response, error) {

	resp, err := h.innerRoundTripper.RoundTrip(req)

	if err != nil || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return resp, err
	}

	// try and read the body to get the formatted error
	var errorResponse ErrorResponse
	var errorResponseRaw []byte
	if resp.Body != http.NoBody {
		errorResponseRaw, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(errorResponseRaw, &errorResponse)
	}

	return resp, HTTPUnsuccessfulError{
		Request:      req,
		Response:     resp,
		ErrorBody:    &errorResponse.Error,
		ErrorBodyRaw: errorResponseRaw,
	}
}
