package http

import (
	"net/http"
	"time"

	"github.com/PennState/go-healthcheck/pkg/health"
)

const (
	statusMeasurementName  = "HTTP/1.1 Status"
	durationMeaurementName = "Latency"
)

var (
	defaultTimeout = 5 * time.Second
)

type Check struct {
	HttpClient http.Client
	// map[name]url ?
	MustPassURLs []string
	MayFailURLs  []string
}

type urlResult struct {
	checks []health.ComponentDetail
	status health.Status
}

func (h Check) Check() ([]health.ComponentDetail, health.Status) {
	var checks []health.ComponentDetail

	mustPassChecks := h.MustPassURLs[:]
	mustPassResults := make(chan urlResult)

	mayFailChecks := h.MayFailURLs[:]
	mayFailResults := make(chan urlResult)

	client := h.HttpClient

	// Ensure the client has a timeout
	if client.Timeout == 0 {
		client = clientWithDefaultTimeout(client)
	}

	for i := range mustPassChecks {
		hc := mustPassChecks[i]
		go checkURL(client, hc, mustPassResults)
	}
	for i := range mayFailChecks {
		hc := mayFailChecks[i]
		go checkURL(client, hc, mayFailResults)
	}

	overallStatus := health.Pass

	for range mustPassChecks {
		urlResult := <-mustPassResults
		if urlResult.status > overallStatus {
			overallStatus = urlResult.status
		}

		checks = append(checks, urlResult.checks...)
	}

	for range mayFailChecks {
		urlResult := <-mayFailResults
		if urlResult.status > overallStatus {
			// MayFailChecks will at most raise the status to Warn if they fail
			overallStatus = health.Warn
		}

		checks = append(checks, urlResult.checks...)
	}

	return checks, overallStatus
}

func checkURL(client http.Client, url string, ch chan urlResult) {
	links := map[string]string{"target": url}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ch <- urlResult{
			checks: []health.ComponentDetail{
				health.ComponentDetail{
					Key: health.Key{
						ComponentName:   url,
						MeasurementName: statusMeasurementName,
					},
					Output:        err.Error(),
					Time:          time.Now().UTC(),
					ComponentType: "component",
					Links:         links,
					Status:        health.Fail,
				}},
			status: health.Fail,
		}
		return
	}

	startTime := time.Now().UTC()
	resp, err := client.Do(req)
	requestDuration := time.Now().UTC().Sub(startTime)

	if err != nil {
		ch <- urlResult{
			checks: []health.ComponentDetail{
				health.ComponentDetail{
					Key: health.Key{
						ComponentName:   url,
						MeasurementName: statusMeasurementName,
					},
					Output:        err.Error(),
					Time:          startTime,
					ComponentType: "component",
					Links:         links,
					Status:        health.Fail,
				}},
			status: health.Fail,
		}
		return
	}

	defer resp.Body.Close()

	status := health.Pass
	if resp.StatusCode/100 != 2 {
		status = health.Fail
	}

	statusCheck := health.ComponentDetail{
		Key: health.Key{
			ComponentName:   url,
			MeasurementName: statusMeasurementName,
		},
		ObservedValue: resp.StatusCode,
		ObservedUnit:  statusMeasurementName,
		Time:          startTime,
		ComponentType: "component",
		Links:         links,
		Status:        status,
	}
	if statusCheck.Status != health.Pass {
		statusCheck.Output = resp.Status
	}

	responseTimeCheck := health.ComponentDetail{
		Key: health.Key{
			ComponentName:   url,
			MeasurementName: durationMeaurementName,
		},
		ObservedValue: requestDuration.String(),
		ObservedUnit:  durationMeaurementName,
		Time:          startTime,
		ComponentType: "component",
		Links:         links,
		Status:        health.Pass,
	}

	ch <- urlResult{
		checks: []health.ComponentDetail{statusCheck, responseTimeCheck},
		status: status,
	}
}

func clientWithDefaultTimeout(client http.Client) http.Client {
	clone := new(http.Client)
	*clone = client
	clone.Timeout = defaultTimeout
	return *clone
}
