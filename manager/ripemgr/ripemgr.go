package ripemgr

//go:generate mockgen -destination mocks.go -package ripemgr . RipeMgr

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type RipeMgr interface {

	// GetProbes return atlas.Probe list filtered by opts map
	GetProbes(opts map[string]string) ([]atlas.Probe, error)

	// CreateMeasurements create ripe atlas.Measurement
	CreateMeasurements(miners []*models.Miner, probeIDs string, t int) ([]*atlas.Measurement, error)

	// GetMeasurementResults get ripe atlas.MeasurementResult from last
	GetMeasurementResults(measurementID int) ([]atlas.MeasurementResult, error)

	// GetMeasurement get ripe atlas.Measurement resource
	GetMeasurement(measurementID int) (*atlas.Measurement, error)
}
