package miners

import (
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerService interface {

	// NewClient connect with Probe api.
	GetMinerIPs(deals []fmgr.VerifiedDeal) []*models.Miner
}
