package miners

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerService interface {

	// ParseMiners parse miners from Filecoin
	ParseMiners(offset uint) []*models.Miner

	// ParseMinersByBlockHeight parse miners from Filecoin for a specific block height
	ParseMinersByBlockHeight(height int64) []*models.Miner

	// GetMiners retrieve miners from Db
	GetMiners() []*models.Miner
}
