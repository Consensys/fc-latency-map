package ripemgr

//go:generate mockgen -destination mocks.go -package ripemgr . RipeMgr

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type RipeMgr interface {

	// GetProbe return a probe by id
	GetProbe(id int) (*atlas.Probe, error)

	// GetProbes return probes list
	GetProbes(opts map[string]string) ([]atlas.Probe, error)

	// GetNearestProbe get the nearest probe from a location
	GetNearestProbe(latitude, longitude float64) (*atlas.Probe, error)

	// CreateMeasurements create ripe measurements
	CreateMeasurements(miners []*models.Miner, probeIDs string, t int) ([]*atlas.Measurement, error)

	// GetMeasurementResults get ripe measurements results from last
	GetMeasurementResults(measurementID, startTime int) ([]atlas.MeasurementResult, error)
}
