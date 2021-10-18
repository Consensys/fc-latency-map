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
	RequestProbes() error

	// ListProbes returns all Probes
	ListProbes() []*models.Probe

	// GetTotalProbes returns probes count
	GetTotalProbes() int64

	// Update handle refresh probes list
	Update() bool

	// ImportProbes from ripe
	ImportProbes()
}
