package measurements

//go:generate mockgen -destination mocks.go -package measurements . MeasurementService

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MeasurementService interface {
	ImportMeasurement(measures []atlas.MeasurementResult)

	GetMiners() []*models.Miner

	GetProbIDs(latitude, longitude float64) []string

	CreateMeasurements([]*atlas.Measurement)

	// getMeasuresLastResultTime load RIPE MeasurementResults
	GetMeasuresLastResultTime() map[int]int
}
