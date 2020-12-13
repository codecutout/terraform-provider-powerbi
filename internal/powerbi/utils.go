package powerbi

import (
	"net/url"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
)

func convertStringToPointer(s string) *string {
	return &s
}

func convertBoolToPointer(b bool) *bool {
	return &b
}

func convertStringSliceToPointer(ss []string) *[]string {
	return &ss
}

func convertToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i := range interfaceSlice {
		stringSlice[i] = interfaceSlice[i].(string)
	}
	return stringSlice
}

func nilIfFalse(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func emptyStringToNil(input string) *string {
	if input == "" {
		return nil
	}
	return &input
}

func isHTTP404Error(err error) bool {
	if httpErr, isHTTPErr := toHTTPUnsuccessfulError(err); isHTTPErr && httpErr.Response.StatusCode == 404 {
		return true
	}
	return false
}

func isHTTP401Error(err error) bool {
	if httpErr, isHTTPErr := toHTTPUnsuccessfulError(err); isHTTPErr && httpErr.Response.StatusCode == 401 {
		return true
	}
	return false
}

func toHTTPUnsuccessfulError(err error) (*powerbiapi.HTTPUnsuccessfulError, bool) {
	if err == nil {
		return nil, false
	}

	if urlErr, isURLErr := err.(*url.Error); isURLErr {
		err = urlErr.Unwrap()
	}

	if httpErr, isHTTPErr := err.(powerbiapi.HTTPUnsuccessfulError); isHTTPErr {
		return &httpErr, true
	}
	return nil, false
}

type wrappedError struct {
	Err          error
	ErrorMessage func(err error) string
}

func (e wrappedError) Error() string {
	return e.ErrorMessage(e.Err)
}
