package measurements

import (
	"log"
	"testing"

	atlas "github.com/keltia/ripe-atlas"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/config"
)

var mgrConfig = config.Config()
var apiKey = mgrConfig.GetString("RIPE_API_KEY")

func TestRipe_GetMeasurementResult(t *testing.T) {
	// t.Skip(true)

	r, err := NewClient(apiKey)
	assert.Nil(t, err)

	// creatred online 32148976
	// created with api - 32221571, 32221572
	// [32221621 32221622]
	got, err := r.GetMeasurementResult(32221571)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.GreaterOrEqual(t, len(got), 5)
	assert.NotNil(t, got[0].Measurement)
	assert.NotNil(t, got[0].Probe)
}

func TestRipe_CreatePing(t *testing.T) {
	// t.Skip(true)

	r, err := NewClient(apiKey)
	assert.Nil(t, err)

	miners := []Miner{
		{Address: "x1234", Ip: []string{
			"213.13.146.142",
			"143.204.98.83",
		}},
	}

	probes := []atlas.ProbeSet{
		{
			Type:      "area",
			Value:     "WW",
			Requested: mgrConfig.GetInt("REQUESTED_PROBES"),
		},
	}

	got, err := r.CreatePing(miners, probes)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.NotNil(t, got.Measurements)
	log.Println(got.Measurements)
}
