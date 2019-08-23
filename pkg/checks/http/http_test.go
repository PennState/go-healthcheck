package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/PennState/go-healthcheck/pkg/health"
	"github.com/stretchr/testify/assert"
)

var (
	responseBody = ioutil.NopCloser(bytes.NewReader([]byte{}))

	successResponse = &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       responseBody,
	}

	notFoundResponse = &http.Response{
		Status:     "404 Not Found",
		StatusCode: http.StatusNotFound,
		Body:       responseBody,
	}

	internalErrorResponse = &http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: http.StatusInternalServerError,
		Body:       responseBody,
	}

	testErr = fmt.Errorf("This is my test error")

	successURL             = "https://httpbin.org/status/200"
	internalServerErrorURL = "https://httpbin.org/status/500"
	notFoundErrorURL       = "https://httpbin.org/status/404"
	longDelayURL           = "https://httpbin.org/delay/10"

	responseMap = responses{
		successURL:             successResponse,
		internalServerErrorURL: internalErrorResponse,
		notFoundErrorURL:       notFoundResponse,
	}

	testClient = http.Client{
		Transport: &testRoundTripper{
			response: responseMap,
		}}

	errClient = http.Client{
		Transport: &testRoundTripper{
			err: testErr,
		}}
)

type responses map[string]*http.Response

type testRoundTripper struct {
	response responses
	err      error
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.err != nil {
		return nil, t.err
	}

	r, ok := t.response[req.URL.String()]
	if !ok {
		return nil, fmt.Errorf("No mocked response for %s", req.URL.String())
	}
	return r, nil
}

func TestNoURLs(t *testing.T) {
	check := Check{
		HttpClient: testClient,
	}

	checks, status := check.Check()

	assert.Equal(t, health.Pass, status)
	assert.Empty(t, checks)
}

func TestSuccessfulMustPass(t *testing.T) {
	check := Check{
		HttpClient:   testClient,
		MustPassURLs: []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Pass, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		assert.Equal(t, health.Pass, check.Status)
	}

	t.Logf("%#v", checks)
}

func TestFailedMustPass(t *testing.T) {
	check := Check{
		HttpClient:   testClient,
		MustPassURLs: []string{internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, health.Fail, check.Status)
			assert.Equal(t, http.StatusInternalServerError, check.ObservedValue)
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, health.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestErrorMustPass(t *testing.T) {
	check := Check{
		HttpClient:   errClient,
		MustPassURLs: []string{longDelayURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 1, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, health.Fail, check.Status)
			assert.Contains(t, check.Output, testErr.Error())
		} else {
			t.Fail()
		}
	}
}

func TestMixedMustPass(t *testing.T) {
	check := Check{
		HttpClient:   testClient,
		MustPassURLs: []string{successURL, internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	passes, failures := 0, 0
	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Status == health.Pass {
				passes++
				assert.Equal(t, http.StatusOK, check.ObservedValue)
			} else {
				failures++
				assert.Equal(t, http.StatusInternalServerError, check.ObservedValue)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, health.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
	assert.Equal(t, 1, passes)
	assert.Equal(t, 1, failures)
}

func TestSuccessfulMayFail(t *testing.T) {
	check := Check{
		HttpClient:  testClient,
		MayFailURLs: []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Pass, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		assert.Equal(t, health.Pass, check.Status)
	}
}

func TestFailedMayFail(t *testing.T) {
	check := Check{
		HttpClient:  testClient,
		MayFailURLs: []string{notFoundErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, health.Fail, check.Status)
			assert.Equal(t, http.StatusNotFound, check.ObservedValue)
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, health.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestErrorMayFail(t *testing.T) {
	check := Check{
		HttpClient:  errClient,
		MayFailURLs: []string{longDelayURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 1, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, health.Fail, check.Status)
			assert.Contains(t, check.Output, testErr.Error())
		} else {
			t.Fail()
		}
	}
}

func TestMixedMayFail(t *testing.T) {
	check := Check{
		HttpClient:  testClient,
		MayFailURLs: []string{successURL, internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	passes, failures := 0, 0
	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Status == health.Pass {
				passes++
				assert.Equal(t, http.StatusOK, check.ObservedValue)
			} else {
				failures++
				assert.Equal(t, http.StatusInternalServerError, check.ObservedValue)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, health.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
	assert.Equal(t, 1, passes)
	assert.Equal(t, 1, failures)
}

func TestAllSuccess(t *testing.T) {
	check := Check{
		HttpClient:   testClient,
		MustPassURLs: []string{successURL},
		MayFailURLs:  []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Pass, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	for _, check := range checks {
		assert.Equal(t, health.Pass, check.Status)
	}
}

func TestWarnRollup(t *testing.T) {
	check := Check{
		HttpClient:   testClient,
		MustPassURLs: []string{successURL},
		MayFailURLs:  []string{internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Key.ComponentName == successURL {
				assert.Equal(t, health.Pass, check.Status)
			} else {
				assert.Equal(t, health.Fail, check.Status)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, health.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestFailRollup(t *testing.T) {
	check := Check{
		HttpClient:   testClient,
		MustPassURLs: []string{notFoundErrorURL},
		MayFailURLs:  []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, health.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Key.ComponentName == successURL {
				assert.Equal(t, health.Pass, check.Status)
			} else {
				assert.Equal(t, health.Fail, check.Status)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, health.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestCloneHttpClient(t *testing.T) {
	client := http.Client{
		Transport: &testRoundTripper{
			response: responseMap,
		},
	}

	assert.Equal(t, time.Duration(0), client.Timeout)

	clone := clientWithDefaultTimeout(client)

	assert.Equal(t, time.Duration(0), client.Timeout)
	assert.Equal(t, defaultTimeout, clone.Timeout)
	assert.Equal(t, client.Transport, clone.Transport)
}
