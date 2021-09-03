package ripeapi



type RipeAPI interface {

	// GetProbes return probes list
	GetProbes() error

	// GetMeasurements return measruments results
	GetMeasurements() error

}