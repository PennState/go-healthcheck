package http

import (
	"net/http"
	"time"

	healthcheck "github.com/PennState/go-healthcheck/pkg/health"
)

const (
	statusMeasurementName  = "HTTP/1.1 Status"
	durationMeaurementName = "Response Time"
)

var (
	defaultTimeout = 5 * time.Second
)

type HTTPCheck struct {
	HttpClient http.Client
	// map[name]url ?
	MustPassURLs []string
	MayFailURLs  []string
}

type urlResult struct {
	checks []healthcheck.Check
	status healthcheck.Status
}

func (h HTTPCheck) Check() ([]healthcheck.Check, healthcheck.Status) {
	var checks []healthcheck.Check

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

	overallStatus := healthcheck.Pass

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
			overallStatus = healthcheck.Warn
		}

		checks = append(checks, urlResult.checks...)
	}

	return checks, overallStatus
}

func checkURL(client http.Client, url string, ch chan urlResult) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ch <- urlResult{
			checks: []healthcheck.Check{
				healthcheck.Check{
					Key: healthcheck.Key{
						ComponentName:   url,
						MeasurementName: statusMeasurementName,
					},
					Output:        err.Error(),
					Time:          time.Now().UTC(),
					ComponentType: "system",
					Links:         map[string]string{"target": url},
					Status:        healthcheck.Fail,
				}},
			status: healthcheck.Fail,
		}
		return
	}

	startTime := time.Now().UTC()
	resp, err := client.Do(req)
	requestDuration := time.Now().UTC().Sub(startTime)

	if err != nil {
		ch <- urlResult{
			checks: []healthcheck.Check{
				healthcheck.Check{
					Key: healthcheck.Key{
						ComponentName:   url,
						MeasurementName: statusMeasurementName,
					},
					Output:        err.Error(),
					Time:          startTime,
					ComponentType: "system",
					Links:         map[string]string{"target": url},
					Status:        healthcheck.Fail,
				}},
			status: healthcheck.Fail,
		}
		return
	}

	defer resp.Body.Close()

	status := healthcheck.Pass
	if resp.StatusCode/100 != 2 {
		status = healthcheck.Fail
	}

	statusCheck := healthcheck.Check{
		Key: healthcheck.Key{
			ComponentName:   url,
			MeasurementName: statusMeasurementName,
		},
		ObservedValue: resp.StatusCode,
		ObservedUnit:  statusMeasurementName,
		Time:          startTime,
		ComponentType: "system",
		Links:         map[string]string{"target": url},
		Status:        status,
	}
	if statusCheck.Status != healthcheck.Pass {
		statusCheck.Output = resp.Status
	}

	responseTimeCheck := healthcheck.Check{
		Key: healthcheck.Key{
			ComponentName:   url,
			MeasurementName: durationMeaurementName,
		},
		ObservedValue: requestDuration.String,
		ObservedUnit:  "Duration",
		Time:          startTime,
		ComponentType: "system",
		Links:         map[string]string{"target": url},
		Status:        healthcheck.Pass,
	}

	ch <- urlResult{
		checks: []healthcheck.Check{statusCheck, responseTimeCheck},
		status: status,
	}
}

func clientWithDefaultTimeout(client http.Client) http.Client {
	clone := new(http.Client)
	*clone = client
	clone.Timeout = defaultTimeout
	return *clone
}
