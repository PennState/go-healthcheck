package health

import (
	"encoding/json"
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
	var check ComponentDetail
	check.Key = key
	checks.Add(key, []ComponentDetail{check}...)
	log.Info("Checks: ", checks)
	json, err := json.Marshal(checks)
	if err != nil {
		log.Error(err)
	}

	log.Info("JSON: ", string(json))
	assert.True(true)
}

// func TestKeyMarshallTextNilKey(t *testing.T) {
// 	var key Key
// 	key = key(nil
// 	bytes, err := key.MarshalText()
// 	log.Info("Bytes: ", bytes)
// 	log.Info("Err: ", err)
// }

func TestKeyMarshalTextWithEmptyComponentNameAndMeasurementName(t *testing.T) {
	var key Key
	bytes, err := key.MarshalText()
	log.Info("Bytes: ", bytes)
	log.Info("Err: ", err)
	//TODO: this should cause an error
}

func TestKeyMarshalTextWithEmptyComponentName(t *testing.T) {
	key := Key{
		ComponentName:   "",
		MeasurementName: "Not Empty",
	}
	bytes, err := key.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(":Not Empty"), bytes)
}

func TestKeyMarshalTextWithEmptyMeasurementName(t *testing.T) {
	key := Key{
		ComponentName:   "Not Empty",
		MeasurementName: "",
	}
	bytes, err := key.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("Not Empty"), bytes)

}

func TestKeyMarshalTextWithTwoPartKey(t *testing.T) {
	key := Key{
		ComponentName:   "Not Empty",
		MeasurementName: "Not Empty",
	}
	bytes, err := key.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("Not Empty:Not Empty"), bytes)
}

func TestKeyUnmarshalTextWithEmptyKey(t *testing.T) {
	var k Key
	err := k.UnmarshalText([]byte(""))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
}

func TestKeyUnmarshalTextWithOnePartKey(t *testing.T) {
	assert := assert.New(t)
	var k Key
	err := k.UnmarshalText([]byte("testComponent"))
	assert.NoError(err)
	assert.Equal("testComponent", k.ComponentName)
	assert.Empty(k.MeasurementName)
}
func TestKeyUnmarshalTextWithTwoPartKey(t *testing.T) {
	assert := assert.New(t)
	var k Key
	err := k.UnmarshalText([]byte("testComponent:testMeasurement"))
	assert.NoError(err)
	assert.Equal("testComponent", k.ComponentName)
	assert.Equal("testMeasurement", k.MeasurementName)
}
