package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	healthcheck "github.com/PennState/go-healthcheck/pkg/health"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRFCExampleCanBeUnmarshaled(t *testing.T) {
	file, err := os.Open("./testdata/rfc.json")
	require.NoError(t, err)
	data, err := ioutil.ReadAll(file)
	require.NoError(t, err)
	var health healthcheck.Health
	err = json.Unmarshal(data, &health)
	require.NoError(t, err)
	log.Info("Health: ", health)
	// TODO: Verify all incoming data is included in result
	// TODO: Compare against "golden file" (or update)
	// TODO: Round-trip the data and compare the source and result JSON
}
