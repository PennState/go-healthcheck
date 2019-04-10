package healthcheck

import (
	"net/http"
)

//State indicates whether the service as-a-whole and the individual checks
//are okay.
type Status string

const (
	Fail Status = "fail"
	Pass Status = "pass"
	Warn Status = "warn"
)

var statusData = []struct {
	Name         string
	ResponseCode int
	Status       Status
}{
	{"Pass", http.StatusOK, Pass},
	{"Fail", http.StatusServiceUnavailable, Fail},
	{"Warn", http.StatusNotFound, Warn},
	{"Undetermined", http.StatusInternalServerError, Fail},
}
