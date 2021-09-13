package miners

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerService interface {

	// ParseMiners parse miners from Filecoin
	ParseMiners() []*models.Miner

	// GetMiners retrieve miners from Db
	GetMiners() []*models.Miner
}
