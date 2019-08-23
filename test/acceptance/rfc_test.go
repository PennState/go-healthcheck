package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/PennState/go-healthcheck/pkg/health"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example() health.Health {
	return health.Health{
		Status:      health.Pass,
		Version:     "1",
		ReleaseId:   "1.2.2",
		Notes:       []string{""},
		Output:      "",
		ServiceId:   "f03e522f-1f44-4062-9b55-9587f91c9c41",
		Description: "health of authz service",
		Checks: health.Checks{
			health.Key{
				ComponentName:   "cassandra",
				MeasurementName: "responseTime",
			}: []health.ComponentDetail{
				health.ComponentDetail{
					ComponentId:   "dfd6cf2b-1b6e-4412-a0b8-f6f7797a60d2",
					ComponentType: "datastore",
					ObservedValue: float64(250),
					ObservedUnit:  "ms",
					Status:        health.Pass,
					AffectedEndpoints: []string{
						"/users/{userId}",
						"/customers/{customerId}/status",
						"/shopping/{anything}",
					},
					Time: timeNoError("2018-01-17T03:36:48Z"),
				},
			},
			health.Key{
				ComponentName:   "cassandra",
				MeasurementName: "connections",
			}: []health.ComponentDetail{
				health.ComponentDetail{
					ComponentId:   "dfd6cf2b-1b6e-4412-a0b8-f6f7797a60d2",
					ComponentType: "datastore",
					ObservedValue: float64(75),
					Status:        health.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
					Links: map[string]string{
						"self": "http://api.example.com/dbnode/dfd6cf2b/health",
					},
				},
			},
			health.Key{
				ComponentName: "uptime",
			}: []health.ComponentDetail{
				health.ComponentDetail{
					ComponentType: "system",
					ObservedValue: float64(1209600.245),
					ObservedUnit:  "s",
					Status:        health.Pass,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
				},
			},
			health.Key{
				ComponentName:   "cpu",
				MeasurementName: "utilization",
			}: []health.ComponentDetail{
				health.ComponentDetail{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					ObservedValue: float64(85),
					ObservedUnit:  "percent",
					Status:        health.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
					AdditionalProperties: map[string]interface{}{
						"node": float64(1),
					},
				},
				health.ComponentDetail{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					ObservedValue: float64(85),
					ObservedUnit:  "percent",
					Status:        health.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
					AdditionalProperties: map[string]interface{}{
						"node": float64(2),
					},
				},
			},
			health.Key{
				ComponentName:   "memory",
				MeasurementName: "utilization",
			}: []health.ComponentDetail{
				health.ComponentDetail{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					ObservedValue: float64(8.5),
					ObservedUnit:  "GiB",
					Status:        health.Warn,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
					AdditionalProperties: map[string]interface{}{
						"node": float64(1),
					},
				},
				health.ComponentDetail{
					ComponentId:   "6fd416e0-8920-410f-9c7b-c479000f7227",
					ComponentType: "system",
					ObservedValue: float64(5500),
					ObservedUnit:  "MiB",
					Status:        health.Pass,
					Time:          timeNoError("2018-01-17T03:36:48Z"),
					AdditionalProperties: map[string]interface{}{
						"node": float64(2),
					},
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
	var health health.Health
	err = json.Unmarshal(data, &health)
	require.NoError(err)
	log.Debug("Health: ", health)

	assert.Equal(Example(), health)

	// TODO: Compare against "golden file" (or update)
	// TODO: Round-trip the data and compare the source and result JSON
}
