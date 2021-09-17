package measurements

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/keltia/ripe-atlas"
)

type MeasurementService interface {
	importMeasurement(measures []atlas.MeasurementResult)

	getMiners() []*models.Miner

	getProbIDs() []string

	createMeasurements([]*atlas.Measurement)

	// getMeasuresLastResultTime load RIPE MeasurementResults
	getMeasuresLastResultTime() map[int]int
}
