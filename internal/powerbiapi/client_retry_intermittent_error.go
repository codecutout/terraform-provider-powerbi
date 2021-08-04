package powerbiapi

import (
	"net/http"
	"time"
)

type retryIntermittentErrorRoundTripper struct {
	innerRoundTripper http.RoundTripper
}

func newRetryIntermittentErrorRoundTripper(next http.RoundTripper) http.RoundTripper {
	return &retryIntermittentErrorRoundTripper{
		innerRoundTripper: next,
	}
}

func (rt *retryIntermittentErrorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	resp, err := rt.innerRoundTripper.RoundTrip(req)

retry:
	for attempts := 1; err == nil && resp.StatusCode == 500 || resp.StatusCode == 400; attempts++ {
		switch attempts {
		case 1:
			// retry immediately. PowerBI API typically responds successfully on a retry
			break
		case 2:
			// gives the service a second
			time.Sleep(1 * time.Second)
			break
		default:
			// we have retried enough
			break retry
		}

		resp, err = rt.innerRoundTripper.RoundTrip(req)
	}

	return resp, err
}
