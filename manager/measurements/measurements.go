package measurements

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type MeasurementService interface {
	RipeCreateMeasurements()

	// RipeCreatePingWithProbes  creates a Ping MeasurementResults for a specific type
	RipeCreatePingWithProbes(miners []*models.Miner, probeIDs string) (*atlas.MeasurementRequest, *atlas.MeasurementResp, error)
	// RipeCreatePing
	RipeCreatePing(miners []*models.Miner, probes []atlas.ProbeSet) (*atlas.MeasurementRequest, *atlas.MeasurementResp, error)

	// RipeGetMeasurement returns MeasurementResults from ID
	RipeGetMeasurement(id int) (*atlas.Measurement, error)

	// GetMeasurementResult
	RipeGetMeasurementResult(id int, start int) ([]atlas.MeasurementResult, error)

	// RipeGetMeasures load RIPE MeasurementResults
	RipeGetMeasures()

	// dbExportData from db to json file
	dbExportData(fn string)

	dbCreate(measurements []*models.Measurement)
}
