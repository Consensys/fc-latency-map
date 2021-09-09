package miners

import (
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerService interface {

	// GetMinerIPs returns miners IP addresses
	GetMinerIPs(deals []fmgr.VerifiedDeal) []*models.Miner
}
