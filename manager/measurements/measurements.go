package measurements

//go:generate mockgen -destination mocks.go -package measurements . MeasurementService

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MeasurementService interface {
	ImportMeasurement(measures []atlas.MeasurementResult)

	GetMiners() []*models.Miner

	GetProbIDs(places []Place, latitude, longitude float64) []string

	UpsertMeasurements([]*atlas.Measurement)

	// GetMeasuresLastResultTime load RIPE MeasurementResults
	GetMeasuresLastResultTime() ([]*models.Measurement, map[int]int)

	PlacesDataSet() ([]Place, error)
	GetMeasurements() []*models.Measurement
}
