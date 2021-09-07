package miners

import (
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerServiceImpl struct {
	DbMgr db.DatabaseMgr
}

func NewMinerServiceImpl(dbMgr db.DatabaseMgr) MinerService {
	return &MinerServiceImpl{
		DbMgr: dbMgr,
	}
}

func (srv *MinerServiceImpl) CreateMiner(miner *models.Miner) error {
	dbMgr := srv.DbMgr
	dbMgr.GetDb().Create(miner)
	return nil
}
