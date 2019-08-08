package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	healthcheck "github.com/PennState/go-healthcheck/pkg/health"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example() healthcheck.Health {
	return healthcheck.Health{
		Status:      healthcheck.Pass,
		Version:     "1",
		ReleaseId:   "1.2.2",
		Notes:       []string{""},
		Output:      "",
		ServiceId:   "f03e522f-1f44-4062-9b55-9587f91c9c41",
		Description: "health of authz service",
		Checks: healthcheck.Checks{
			healthcheck.Key{
				ComponentName:   "cassandra",
				MeasurementName: "responseTime",
			}: []healthcheck.Check{
				healthcheck.Check{
					ComponentId:   "dfd6cf2b-1b6e-4412-a0b8-f6f7797a60d2",
					ComponentType: "datastore",
					ObservedValue: float64(250),
					ObservedUnit:  "ms",
					Status:        healthcheck.Pass,
					AffectedEndpoints: []string{
						"/users/{userId}",
						"/customers/{customerId}/status",
						"/shopping/{anything}",
					},
					Time: timeNoError("2018-01-17T03:36:48Z"),
				},
			},
			healthcheck.Key{
				ComponentName:   "cassandra",
				MeasurementName: "connections",
			}: []healthcheck.Check{
				healthcheck.Check{
					ComponentId:   "dfd6cf2b-1b6e-4412-a0b8-f6f7797a60d2",
					ComponentType: "datastore",
					ObservedValue: float64(75),
					Status:        healthcheck.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
					Links: map[string]string{
						"self": "http://api.example.com/dbnode/dfd6cf2b/health",
					},
				},
			},
			healthcheck.Key{
				ComponentName: "uptime",
			}: []healthcheck.Check{
				healthcheck.Check{
					ComponentType: "system",
					ObservedValue: float64(1209600.245),
					ObservedUnit:  "s",
					Status:        healthcheck.Pass,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
				},
			},
			healthcheck.Key{
				ComponentName:   "cpu",
				MeasurementName: "utilization",
			}: []healthcheck.Check{
				healthcheck.Check{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					// TODO: Missing the Node field
					ObservedValue: float64(85),
					ObservedUnit:  "percent",
					Status:        healthcheck.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
				},
				healthcheck.Check{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					// TODO: Missing the Node field
					ObservedValue: float64(85),
					ObservedUnit:  "percent",
					Status:        healthcheck.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
				},
			},
			healthcheck.Key{
				ComponentName:   "memory",
				MeasurementName: "utilization",
			}: []healthcheck.Check{
				healthcheck.Check{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					// TODO: Missing the Node field
					ObservedValue: float64(8.5),
					ObservedUnit:  "GiB",
					Status:        healthcheck.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
				},
				healthcheck.Check{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					// TODO: Missing the Node field
					ObservedValue: float64(5500),
					ObservedUnit:  "MiB",
					Status:        healthcheck.Pass,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
				},
			},
		},
		Links: map[string]string{
			"about":                          "http://api.example.com/about/authz",
			"http://api.x.io/rel/thresholds": "http://api.x.io/about/authz/thresholds",
		},
	}
}

func timeNoError(t string) time.Time {
	time, _ := time.Parse(time.RFC3339, t)
	return time
}

func TestRFCExampleCanBeUnmarshaled(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	file, err := os.Open("./testdata/rfc.json")
	require.NoError(err)
	data, err := ioutil.ReadAll(file)
	require.NoError(err)
	var health healthcheck.Health
	err = json.Unmarshal(data, &health)
	require.NoError(err)
	log.Debug("Health: ", health)

	// TODO: Add missing "node" fields to CPU and memory utilization
	assert.Equal(Example(), health)

	// TODO: Compare against "golden file" (or update)
	// TODO: Round-trip the data and compare the source and result JSON
}
