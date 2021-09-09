package miners

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/util"
)

type MinerServiceImpl struct {
	DbMgr db.DatabaseMgr
	FMgr  fmgr.FilecoinMgr
}

func NewMinerServiceImpl(dbMgr db.DatabaseMgr, fMgr fmgr.FilecoinMgr) MinerService {
	return &MinerServiceImpl{
		DbMgr: dbMgr,
		FMgr:  fMgr,
	}
}

func (srv *MinerServiceImpl) GetMinerIPs(deals []fmgr.VerifiedDeal) []*models.Miner {
	var miners = []*models.Miner{}
	for _, deal := range deals {
		provider := deal.Provider
		minerInfo, err := srv.FMgr.GetMinerInfo(provider)
		if err != nil {
			log.Printf("unable to get miner info: %s", provider)
			continue
		}
		ips := util.IpAddress(util.MultiAddrs(minerInfo.Multiaddrs))
		ip := ""
		if len(ips) > 0 {
			ip = ips[0]
		}
		miner := &models.Miner{
			Address: deal.Provider.String(),
			Ip:      ip,
		}
		miners = append(miners, miner)
		srv.DbMgr.GetDb().Create(miner)
	}
	return miners
}
