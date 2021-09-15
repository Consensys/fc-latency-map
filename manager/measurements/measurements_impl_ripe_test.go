package measurements

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/config"
)

func TestName(t *testing.T) {
	t.Skip(true)
	cfg := atlas.Config{
		APIKey: "",
	}
	client, err := atlas.NewClient(cfg)
	if err != nil {
		return

	}
	err = client.DeleteMeasurement(32148976)
	if err != nil {
		fmt.Print(err)
		return
	}
}

func TestRipe_GetMeasurementResult(t *testing.T) {
	// t.Skip(true)

	r := NewHandler()

	// creatred online 32148976
	// created with api - 32221571, 32221572
	// [32221621 32221622]
	got, err := (*r.Service).RipeGetMeasurementResult(32221571, 0)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.GreaterOrEqual(t, len(got), 5)
	assert.NotNil(t, got)
}

func TestRipe_CreatePing(t *testing.T) {
	// t.Skip(true)

	r := NewHandler()

	miners := []*models.Miner{
		{Address: "x1234", Ip: "213.13.146.142,143.204.98.83"},
	}
	mgrConfig := config.NewConfig()

	probes := []atlas.ProbeSet{
		{
			Type:      "area",
			Value:     "WW",
			Requested: mgrConfig.GetInt("RIPE_REQUESTED_PROBES"),
		},
	}

	_, got, err := (*r.Service).RipeCreatePing(miners, probes)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.NotNil(t, got.Measurements)
	log.Println(got.Measurements)
}

func TestRipe_CreatePingWithProbID(t *testing.T) {
	// t.Skip(true)

	r := NewHandler()

	miners := []*models.Miner{
		{Address: "f0883203", Ip: "10.6.13.218"},
		{Address: "f0883203", Ip: "213.13.146.142,143.204.98.83"},
		// {Address: "xminer20210910", Ip: "213.13.146.142,143.204.98.83"},
		// {Address: "xminer20210911", Ip: "213.13.146.142,143.204.98.83"},
	}

	_, got, err := (*r.Service).RipeCreatePingWithProbes(miners, "1001065,6252")

	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.NotNil(t, got.Measurements)
	log.Println(got.Measurements)
}

func TestRipe_GetMeasures(t *testing.T) {
	// t.Skip(true)
	r := NewHandler()
	r.GetMeasures()
}
