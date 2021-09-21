package ripemgr

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type RipeMgr interface {

	// GetProbe return a probe by id
	GetProbe(id int) (*atlas.Probe, error)

	// GetProbes return probes list
	GetProbes(opts map[string]string) ([]atlas.Probe, error)

	// GetNearestProbes get the nearest probe from a location
	GetNearestProbe(latitude, longitude string) (*atlas.Probe, error)

	// CreateMeasurements create ripe measurements
	CreateMeasurements(miners []*models.Miner, probeIDs string) ([]*atlas.Measurement, error)

	// GetMeasurementResults get ripe measurements results from last
	GetMeasurementResults(measures map[int]int) ([]atlas.MeasurementResult, error)
}
