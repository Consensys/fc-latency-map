package measurements

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type MeasurementService interface {

	// GetMeasurement returns Measurement from ID
	GetRipeMeasurement(id int) (*atlas.Measurement, error)

	// CreatePing creates a Ping Measurement
	CreatePing(miners []*models.Miner, probes []atlas.ProbeSet) (*atlas.MeasurementResp, error)

	// CreatePingByType  creates a Ping Measurement for a specific type
	CreatePingProbes(miners []*models.Miner, probeType, value string) (*atlas.MeasurementResp, error)

	// GetMeasurementResult
	GetRipeMeasurementResult(id int) ([]atlas.MeasurementResult, error)

	// GetMeasures load RIPE Measurement
	GetRipeMeasures(tag string)

	// ExportData from db to json file
	ExportDbData(fn string)

	CreateMeasurements()
}
