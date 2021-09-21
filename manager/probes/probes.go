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

	// RequestProbes returns Probes from Ripe
	RequestProbes() ([]atlas.Probe, error)

	// GetAllProbes returns all Probes
	GetAllProbes() []*models.Probe

	// Update handle refresh probes list
	Update()
}
