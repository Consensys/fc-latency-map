package measurements

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type MeasurementService interface {
	RipeCreateMeasurements()

	// CreatePingByType  creates a Ping MeasurementResults for a specific type
	RipeCreatePingWithProbes(miners []*models.Miner, probeIDs string) (*atlas.MeasurementRequest, *atlas.MeasurementResp, error)

	RipeCreatePing(miners []*models.Miner, probes []atlas.ProbeSet) (*atlas.MeasurementRequest, *atlas.MeasurementResp, error)

	// GetMeasurement returns MeasurementResults from ID
	RipeGetMeasurement(id int) (*atlas.Measurement, error)

	// GetMeasurementResult
	RipeGetMeasurementResult(id int, start int) ([]atlas.MeasurementResult, error)

	// GetMeasures load RIPE MeasurementResults
	RipeGetMeasures()

	// ExportData from db to json file
	dbExportData(fn string)

	dbCreate(measurements []*models.Measurement)
}
