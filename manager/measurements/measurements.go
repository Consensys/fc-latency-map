package measurements

//go:generate mockgen -destination mocks.go -package measurements . MeasurementService

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MeasurementService interface {
	ImportMeasurement(measures []atlas.MeasurementResult)

	GetMinersWithGeolocation() []*models.Miner

	GetProbIDs(places []Place, latitude, longitude float64) []string

	UpsertMeasurements([]*atlas.Measurement)

	// GetMeasuresLastResultTime load RIPE MeasurementResults
	GetMeasuresLastResultTime() ([]*models.Measurement, map[int]int)

	// GetLocationsAsPlaces returns slice of measurements.Place
	GetLocationsAsPlaces() ([]Place, error)

	// GetMeasurementsRunning models.Measurement running on ripe
	GetMeasurementsRunning() []*models.Measurement
}
