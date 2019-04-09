package health

import (
	"net/http"
)

//State indicates whether the service as-a-whole and the individual checks
//are okay.
type State string

const (
	Down State = "DOWN" //Down indicates the service is down
	Up   State = "UP"   //Up indicates the service is up
)

var statusData = []struct {
	Name         string
	ResponseCode int
	Status       State
}{
	{"Ok", http.StatusOK, Up},
	{"Error", http.StatusServiceUnavailable, Down},
	{"Undetermined", http.StatusInternalServerError, Down},
}
