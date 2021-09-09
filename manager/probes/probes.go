package probes

import (
	atlas "github.com/keltia/ripe-atlas"
)

type Miner struct {
	Address string
	Ip      []string
}

type API interface {

	// GetProbe returns Probe from ID
	GetProbe(id int) (m *atlas.Probe, err error)

	// GetProbes returns Probes by country code
	GetProbes(countryCode string) ([]atlas.Probe, error)

	// GetBestProbes returns best Probes from Probes list
	GetBestProbes(countryProbes []atlas.Probe) (atlas.Probe, error)

	// GetAllProbes returns all Probes
	GetAllProbes() ([]atlas.Probe, error)
}
