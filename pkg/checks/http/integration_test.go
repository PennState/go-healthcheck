// +build integration
// run with go test --tags=integration

package http

import (
	"net/http"
	"testing"

	healthcheck "github.com/PennState/go-healthcheck/pkg/health"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulMustPassIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Pass, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		assert.Equal(t, healthcheck.Pass, check.Status)
	}
}

func TestFailedMustPassIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, healthcheck.Fail, check.Status)
			assert.Equal(t, http.StatusInternalServerError, check.ObservedValue)
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, healthcheck.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestErrorMustPassIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{longDelayURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 1, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, healthcheck.Fail, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestMixedMustPassIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{successURL, internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	passes, failures := 0, 0
	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Status == healthcheck.Pass {
				passes++
				assert.Equal(t, http.StatusOK, check.ObservedValue)
			} else {
				failures++
				assert.Equal(t, http.StatusInternalServerError, check.ObservedValue)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, healthcheck.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
	assert.Equal(t, 1, passes)
	assert.Equal(t, 1, failures)
}

func TestSuccessfulMayFailIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:  http.Client{},
		MayFailURLs: []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Pass, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		assert.Equal(t, healthcheck.Pass, check.Status)
	}
}

func TestFailedMayFailIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:  http.Client{},
		MayFailURLs: []string{notFoundErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 2, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, healthcheck.Fail, check.Status)
			assert.Equal(t, http.StatusNotFound, check.ObservedValue)
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, healthcheck.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestErrorMayFailIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:  errClient,
		MayFailURLs: []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 1, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			assert.Equal(t, healthcheck.Fail, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestMixedMayFailIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:  http.Client{},
		MayFailURLs: []string{successURL, internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	passes, failures := 0, 0
	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Status == healthcheck.Pass {
				passes++
				assert.Equal(t, http.StatusOK, check.ObservedValue)
			} else {
				failures++
				assert.Equal(t, http.StatusInternalServerError, check.ObservedValue)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, healthcheck.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
	assert.Equal(t, 1, passes)
	assert.Equal(t, 1, failures)
}

func TestAllSuccessIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{successURL},
		MayFailURLs:  []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Pass, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	for _, check := range checks {
		assert.Equal(t, healthcheck.Pass, check.Status)
	}
}

func TestWarnRollupIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{successURL},
		MayFailURLs:  []string{internalServerErrorURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Warn, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Key.ComponentName == successURL {
				assert.Equal(t, healthcheck.Pass, check.Status)
			} else {
				assert.Equal(t, healthcheck.Fail, check.Status)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, healthcheck.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}

func TestFailRollupIntegration(t *testing.T) {
	check := HTTPCheck{
		HttpClient:   http.Client{},
		MustPassURLs: []string{notFoundErrorURL},
		MayFailURLs:  []string{successURL},
	}

	checks, status := check.Check()

	assert.Equal(t, healthcheck.Fail, status)
	assert.NotEmpty(t, checks)
	assert.Equal(t, 4, len(checks))

	for _, check := range checks {
		if check.Key.MeasurementName == statusMeasurementName {
			if check.Key.ComponentName == successURL {
				assert.Equal(t, healthcheck.Pass, check.Status)
			} else {
				assert.Equal(t, healthcheck.Fail, check.Status)
			}
		} else if check.Key.MeasurementName == durationMeaurementName {
			assert.Equal(t, healthcheck.Pass, check.Status)
		} else {
			t.Fail()
		}
	}
}
