package miners

import "github.com/ConsenSys/fc-latency-map/manager/models"

type MinerService interface {

	// NewClient connect with Probe api.
	CreateMiner(miner *models.Miner) error
}
