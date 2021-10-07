package miners

//go:generate mockgen -destination mocks.go -package miners . MinerService

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerService interface {

	// GetAllMiners GetMinersWithGeoLocation retrieve miners from db
	GetAllMiners() []*models.Miner

	// ParseMinersByBlockOffset parse miners from Filecoin
	ParseMinersByBlockOffset(offset int) []*models.Miner

	// ParseMinersByBlockHeight parse miners from Filecoin for a specific block height
	ParseMinersByBlockHeight(height int64) []*models.Miner
}
