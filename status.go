package healthcheck

import (
	"net/http"
)

//State indicates whether the service as-a-whole and the individual checks
//are okay.
type Status string

const (
	StatusFail string = "fail"
	StatusPass string = "pass"
	StatusWarn string = "warn"
)

var statusData = []struct {
	Name         string
	ResponseCode int
	Status       State
}{
	{"Pass", http.StatusOK, StatusPass},
	{"Fail", http.StatusServiceUnavailable, StatusFail},
	{"Warn", nil, StatusWarn},
	{"Undetermined", http.StatusInternalServerError, Down},
}
