package measurements

//go:generate mockgen -destination mocks.go -package measurements . MeasurementService

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MeasurementService interface {
	ImportMeasurement(measures []atlas.MeasurementResult)

	GetMinersWithGeolocation() []*models.Miner

	getProbIDs(places []Place, latitude, longitude float64) []string

	UpsertMeasurements([]*atlas.Measurement)

	// GetLocationsAsPlaces returns slice of measurements.Place
	getLocationsAsPlaces() ([]Place, error)

	// GetMeasurementsRunning models.Measurement running on ripe
	getMeasurementsRunning() []*models.Measurement
}
