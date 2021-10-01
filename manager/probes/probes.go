package probes

//go:generate mockgen -destination mocks.go -package probes . ProbeService

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type Miner struct {
	Address string
	IP      []string
}

type ProbeService interface {

	// RequestProbes returns Probes from Ripe
	RequestProbes() ([]*models.Probe, error)

	// GetAllProbes returns all Probes
	GetAllProbes() []*models.Probe

	// Update handle refresh probes list
	Update()
}
