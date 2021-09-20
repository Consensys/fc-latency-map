package ripemgr

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type RipeService interface {

	// createMeasurements create ripe measurements
	createMeasurements(miners []*models.Miner, probeIDs string) ([]*atlas.Measurement, error)

	// getMeasurementResults get ripe measurements results from last
	getMeasurementResults(measures map[int]int) ([]atlas.MeasurementResult, error)
}
