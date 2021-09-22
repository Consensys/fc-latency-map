package measurements

//go:generate mockgen -destination mocks/measurements_impl.go -package measurements . MeasurementService

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MeasurementService interface {
	importMeasurement(measures []atlas.MeasurementResult)

	getMiners() []*models.Miner

	getProbIDs() []string

	createMeasurements([]*atlas.Measurement)

	// getMeasuresLastResultTime load RIPE MeasurementResults
	getMeasuresLastResultTime() map[int]int
}
