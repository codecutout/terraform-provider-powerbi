package powerbiapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type retryTooManyRequestsRoundTripper struct {
	innerRoundTripper http.RoundTripper
}

func newRetryTooManyRequestsRoundTripper(next http.RoundTripper) http.RoundTripper {
	return &retryTooManyRequestsRoundTripper{
		innerRoundTripper: next,
	}
}

func (rt *retryTooManyRequestsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	resp, err := rt.innerRoundTripper.RoundTrip(req)

retry:
	for attempts := 1; err == nil && resp.StatusCode == 429; attempts++ {
		switch attempts {
		case 1:
			// retry immediately. PowerBI API typically responds successfully on a retry
			break
		case 2:
			// if failed again then we will respect the retry-after header
			// these can unfortunately be anywhere up to a minute
			time.Sleep(readRetryAfter(resp, 5*time.Second))
		case 3:
			// respect retry-after again
			time.Sleep(readRetryAfter(resp, 10*time.Second))
		default:
			// we have retried enough
			break retry
		}

		resp, err = rt.innerRoundTripper.RoundTrip(req)
	}

	return resp, err
}

func readRetryAfter(resp *http.Response, fallback time.Duration) time.Duration {
	waitSeconds, parseErr := extractHeaderAsInteger(resp, "Retry-After")
	if parseErr != nil {
		return fallback
	}
	return time.Duration(waitSeconds+1) * time.Second
}

func extractHeaderAsInteger(resp *http.Response, key string) (int, error) {
	headerValues, ok := resp.Header[key]
	if !ok || len(headerValues) == 0 {
		return 0, fmt.Errorf("Response does not contain key '%s'", key)
	}
	return strconv.Atoi(headerValues[0])
}
