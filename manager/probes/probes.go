package probes

import (
	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type Miner struct {
	Address string
	Ip      []string
}

type ProbeService interface {

	// GetProbe returns Probe from ID
	GetProbe(id int) (m *atlas.Probe, err error)

	// GetProbes returns Probes by country code
	GetProbes(countryCode string) ([]atlas.Probe, error)

	// GetBestProbes returns best Probes from Probes list
	GetBestProbes(countryProbes []atlas.Probe) (atlas.Probe, error)

	// RequestAllProbes returns all Probes from Ripe
	RequestAllProbes() ([]atlas.Probe, error)
	
	// GetAllProbes returns all Probes
	GetAllProbes() []*models.Probe
	
	// Update handle refresh probes list
	Update()
}
