package healthcheck

import (
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

type Checker interface {
	Check() ([]Check, Status)
}

//Check is an individual element of the array returned for a given Checks
//key.
//
//See: https://inadarei.github.io/rfc-healthcheck/#the-checks-object
type Check struct {
	Key               Key               `json:"-"`
	ComponentId       string            `json:",omitempty"`
	ComponentType     string            `json:",omitempty"`
	ObservedValue     interface{}       `json:",omitempty"`
	ObservedUnit      string            `json:",omitempty"`
	Status            Status            `json:",omitempty"`
	AffectedEndpoints []string          `json:",omitempty"`
	Time              time.Time         `json:",omitempty"`
	Output            string            `json:",omitempty"`
	Links             map[string]string `json:",omitempty"`
}

//Key provides a composite key denoting the component name and measurement
//name of a health checks.
//
//See: https://inadarei.github.io/rfc-healthcheck/#the-checks-object
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
// Checks
//

type Checks map[Key][]Check

// AddChecks is a convenience method that adds a result for a []
// whether or not its key was previously known.
func (c Checks) AddChecks(check ...Check) {
	for _, v := range check {
		k := v.Key
		l, exists := c[k]
		if !exists {
			l = make([]Check, 1)
		}
		l = append(l, v)
		c[k] = l
	}
}
