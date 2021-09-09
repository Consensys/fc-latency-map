package probes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/config"
)

var mgrConfig = config.NewConfig()
var apiKey = mgrConfig.GetString("RIPE_API_KEY")

func Test_GetAllProbes(t *testing.T) {
	// t.Skip(true)

	r, err := NewClient(apiKey)
	assert.Nil(t, err)

	probes, err := r.GetAllProbes()
	assert.Nil(t, err)

	fmt.Printf("Get %v probes", len(probes))
	assert.GreaterOrEqual(t, len(probes), 1)
}
