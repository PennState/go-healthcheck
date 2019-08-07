package healthcheck

import (
	"fmt"
	"net/http"
	"strings"
)

//State indicates whether the service as-a-whole and the individual checks
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
