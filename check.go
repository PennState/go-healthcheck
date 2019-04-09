package healthcheck

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

type Checker interface {
	Check() Check
}

type Check struct {
	Name   string
	Status State
	Data   map[string]interface{}
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
	cn, err := state.Token(true, func(r rune) bool {
		return r != ':'
	})
	if err != nil {
		return err
	}
	k.ComponentName = string(cn)

	_, _, err = state.ReadRune()
	if err != nil && err == io.ErrUnexpectedEOF {
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

type Checks map[Key][]Check

func (c Checks) MarshalJSON() {
	for k, v := range c {
		log.Info("Key: ", k, ", Value: ", v)
	}
}

func (c *Checks) UnmarshalJSON() {

}
