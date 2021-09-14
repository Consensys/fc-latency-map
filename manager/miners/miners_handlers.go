package miners

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/spf13/viper"
)

type MinerHandler struct {
	Conf *viper.Viper
	MSer *MinerService
}

func NewMinerHandler() *MinerHandler {
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
	mSer := NewMinerServiceImpl(conf, &dbMgr, &fMgr)
	return &MinerHandler{
		Conf: conf,
		MSer: &mSer,
	}
}

func (mHdl *MinerHandler) MinersUpdate(offset string) {
	if len(strings.TrimSpace(offset)) == 0 {
		off := mHdl.Conf.GetUint("FILECOIN_BLOCKS_OFFSET")
		(*mHdl.MSer).ParseMiners(off)
	} else {
		off, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			fmt.Println("Error: provided offset is not a valid integer")
			return
		}
		(*mHdl.MSer).ParseMiners(uint(off))
	}
}

func (mHdl *MinerHandler) MinersParse(height int64) {
	(*mHdl.MSer).ParseMinersByBlockHeight(height)
}
