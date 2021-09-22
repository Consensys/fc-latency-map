package miners

//go:generate mockgen -destination mocks/mocks.go -package miners . MinerService

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerService interface {

	// GetMiners retrieve miners from db
	GetAllMiners() []*models.Miner

	// ParseMiners parse miners from Filecoin
	ParseMiners(offset uint) []*models.Miner

	// ParseMinersByBlockHeight parse miners from Filecoin for a specific block height
	ParseMinersByBlockHeight(height int64) []*models.Miner
}
