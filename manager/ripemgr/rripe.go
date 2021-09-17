package ripemgr

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type Service interface {
	createMeasurements(miners []*models.Miner, probeIDs string) ([]*atlas.Measurement, error)

	getMeasurementResults(measures map[int]int) ([]atlas.MeasurementResult, error)
}
