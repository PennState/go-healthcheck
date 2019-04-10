package healthcheck

import (
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

type Checker interface {
	Check() Check
}

//Check is an individual element of the array returned for a given Checks
//key.
//
//See: https://inadarei.github.io/rfc-healthcheck/#the-checks-object
type Check struct {
	ComponentId       string
	ComponentType     string
	ObservedValue     interface{}
	ObservedUnit      string
	Status            Status
	AffectedEndpoints []string
	Time              time.Time
	Output            string
	Links             []string
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

func (c Checks) MarshalJSON() ([]byte, error) {
	log.Info("Got here!")
	for k, v := range c {
		log.Info("Key: ", k, ", Value: ", v)
	}
	log.Info("Got here too!")
	return []byte("{}"), nil
}

func (c *Checks) UnmarshalJSON(json []byte) error {
	return nil
}

func (c Checks) AddCheck(key Key, check Check) {
	l, exists := c[key]
	if !exists {
		l = make([]Check, 1)
	}
	l = append(l, check)
	c[key] = l
}
