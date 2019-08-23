//go:generate additional-properties

package health

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//
// Health - https://inadarei.github.io/rfc-healthcheck/#api-health-response
//

type Health struct {
	Status      Status            `json:"status"`
	Version     string            `json:"version,omitempty"`
	ReleaseId   string            `json:"releaseId,omitempty"`
	Notes       []string          `json:"notes,omitempty"`
	Output      string            `json:"output,omitempty"`
	Checks      Checks            `json:"checks,omitempty"`
	Links       map[string]string `json:"links,omitempty"`
	ServiceId   string            `json:"serviceId,omitempty"`
	Description string            `json:"description,omitempty"`
}

//
// Status - https://inadarei.github.io/rfc-healthcheck/#status
//

//Status indicates whether the service as-a-whole and the individual checks
//are okay.
type Status int

//These must be organized from least to greatest severity and must be in
//the same order as the data structure below. (Replace with go-enumeration
//when it's available.)
const (
	Pass Status = iota
	Warn Status = 1
	Fail Status = 2
)

var statusData = []struct {
	name         string
	responseCode int
}{
	{"Pass", http.StatusOK},
	{"Warn", http.StatusOK},
	{"Fail", http.StatusServiceUnavailable},
	//{"Undetermined", http.StatusOK},
}

func ParseStatus(input string) (Status, error) {
	switch strings.ToLower(input) {
	case strings.ToLower(Pass.String()):
		return Pass, nil
	case strings.ToLower(Fail.String()):
		return Fail, nil
	case strings.ToLower(Warn.String()):
		return Warn, nil
	}
	return Pass, fmt.Errorf("Couldn't parse Status with value: %v", input)
}

func (s Status) MarshalText() ([]byte, error) {
	return []byte(strings.ToLower(s.String())), nil
}

func (s Status) Max(other Status) Status {
	if s.Severity() >= other.Severity() {
		return s
	}
	return other
}

func (s Status) Severity() int {
	return int(s)
}

func (s Status) String() string {
	return statusData[s].name
}

func (s Status) StatusCode() int {
	return statusData[s].responseCode
}

func (s *Status) UnmarshalText(json []byte) error {
	st, err := ParseStatus(string(json))
	if err == nil {
		*s = st
	}
	return err
}

//
// Checks - https://inadarei.github.io/rfc-healthcheck/#the-checks-object
//

type Checks map[Key][]ComponentDetail

// Add is a convenience method that adds one or more ComponentDetail
// objects into the Checks receiver for the provided key (whether or
// not that key was previously known).
func (c Checks) Add(key Key, details ...ComponentDetail) {
	v, ok := c[key]
	if !ok {
		v = []ComponentDetail{}
	}
	c[key] = append(v, details...)
}

// Merge is a convenience method that adds more or more Checks objects
// into the Checks receiver.
func (c Checks) Merge(checks ...Checks) {
	for _, check := range checks {
		for k, v := range check {
			c.Add(k, v...)
		}
	}
}

//
// Key - https://inadarei.github.io/rfc-healthcheck/#the-checks-object
//

//Key provides a composite key denoting the component name and measurement
//name of a health checks.
type Key struct {
	ComponentName   string
	MeasurementName string
}

func (k Key) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k *Key) Scan(state fmt.ScanState, verb rune) error {
	_, _, err := state.ReadRune()
	if err != nil {
		log.Debug("Missing componentName but it's mandatory")
		return err
	}
	state.UnreadRune()

	cn, err := state.Token(true, func(r rune) bool {
		return r != ':'
	})
	if err != nil {
		return err
	}
	k.ComponentName = string(cn)

	_, _, err = state.ReadRune()
	if err != nil && (err == io.ErrUnexpectedEOF || err == io.EOF) {
		log.Debug("There was no separator (:) found so there is no measurementName")
		return nil
	}
	if err != nil {
		return err
	}

	mn, err := state.Token(true, func(r rune) bool {
		return true
	})
	if err != nil {
		return err
	}
	k.MeasurementName = string(mn)

	return nil
}

func (k Key) String() string {
	if k.MeasurementName == "" {
		return k.ComponentName
	}
	return k.ComponentName + ":" + k.MeasurementName
}

func (k *Key) UnmarshalText(text []byte) error {
	_, err := fmt.Sscan(string(text), k)
	return err
}

//
// ComponentDetail - https://inadarei.github.io/rfc-healthcheck/#the-checks-object
//

// ComponentDetail is an individual element of the array returned for a
// given Checks key.
type ComponentDetail struct {
	Key                  Key                    `json:"-"`
	ComponentId          string                 `json:"componentId,omitempty"`
	ComponentType        string                 `json:"componentType,omitempty"`
	ObservedValue        interface{}            `json:"observedValue,omitempty"`
	ObservedUnit         string                 `json:"observedUnit,omitempty"`
	Status               Status                 `json:"status,omitempty"`
	AffectedEndpoints    []string               `json:"affectedEndpoints,omitempty"`
	Time                 time.Time              `json:"time,omitempty"`
	Output               string                 `json:"output,omitempty"`
	Links                map[string]string      `json:"links,omitempty"`
	AdditionalProperties map[string]interface{} `json:"*,omitempty"`
}
