package miners

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
)

type MinerHandler struct {
	MSer MinerService
}

func NewMinerHandler() MinerHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	nodeUrl := conf.GetString("FILECOIN_NODE_URL")
	fMgr, err := fmgr.NewFilecoinImpl(nodeUrl)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	mSer := NewMinerServiceImpl(conf, dbMgr, fMgr)
	return MinerHandler{
		MSer: mSer,
	}
}

func (mHdl *MinerHandler) MinersUpdate() {
	mHdl.MSer.ParseMiners()
}
