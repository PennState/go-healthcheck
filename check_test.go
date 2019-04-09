package healthcheck

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	assert := assert.New(t)
	checks := make(Checks)
	key := Key{
		ComponentName:   "testComponent",
		MeasurementName: "testMeasurement",
	}
	var check Check
	var checkSlice []Check
	checkSlice = append(checkSlice, check)
	log.Info("Check: ", check)
	checks[key] = checkSlice
	log.Info("Checks: ", checks)
	json, err := json.Marshal(checks)
	if err != nil {
		log.Error(err)
	}

	log.Info("JSON: ", string(json))
	assert.True(false)
}

func TestKeyScanWithEmptyKey(t *testing.T) {
	var k Key
	n, err := fmt.Sscan("", &k)
	log.Info("Size: ", n)
	log.Info("Error: ", err)
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
}

func TestKeyScanWithOnePartKey(t *testing.T) {
	var k Key
	n, err := fmt.Sscan("testComponent", &k)
	log.Info("Size: ", n)
	log.Info("Error: ", err)
	assert.Equal(t, "testComponent", k.ComponentName)
	assert.Empty(t, k.MeasurementName)
}
func TestKeyScanWithTwoPartKey(t *testing.T) {
	var k Key
	fmt.Sscan("testComponent:testMeasurement", &k)
	assert.Equal(t, "testComponent", k.ComponentName)
	assert.Equal(t, "testMeasurement", k.MeasurementName)
}
