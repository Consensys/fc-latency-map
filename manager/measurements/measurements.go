package measurements

import (
    atlas "github.com/keltia/ripe-atlas"
)

type Miner struct {
    Address string
    Ip []string
}

type MeasurementResults struct {
    Measurement atlas.MeasurementResult
    Probe atlas.Probe
}

type API interface {

    // NewClient connect with Measurement api.
    NewClient(apiToken string, cfgs ...atlas.Config) error

    // GetMeasurement returns Measurement from ID
    GetMeasurement(id int) (m *atlas.Measurement, err error)

    // CreatePing creates a Ping Measurement
    CreatePing(miners []Miner) (m *atlas.MeasurementResp, err error)

    // GetMeasurementResult returns all the probe Measurements for the MeasurementID
    GetMeasurementResult(id int) (m *MeasurementResults, err error)
}
