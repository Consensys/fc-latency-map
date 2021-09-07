package probes

import (
	"fmt"
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/stretchr/testify/assert"
)

var mgrConfig = config.Config()
var apiKey string = mgrConfig.GetString("RIPE_API_KEY")

func Test_GetAllProbes(t *testing.T) {
    // t.Skip(true)

    r := &Ripe{}
    err := r.NewClient(apiKey)
    assert.Nil(t, err)

    probes, err := r.GetAllProbes()
		assert.Nil(t, err)

		fmt.Printf("Get %v probes", len(probes))
		assert.GreaterOrEqual(t, len(probes), 1)
}