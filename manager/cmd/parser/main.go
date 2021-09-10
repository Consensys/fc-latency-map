package main

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/parser"
)

var mgrConfig = config.NewConfig()
var nodeUrl string = mgrConfig.GetString("FILECOIN_NODE_URL")

func main() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	fMgr, err := filecoinmgr.NewFilecoinImpl(nodeUrl)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	mSer := miners.NewMinerServiceImpl(conf, dbMgr, fMgr)
	parser.Parse(mSer)

	var miners = mSer.GetMiners()
	for _, miner := range miners {
		fmt.Printf("Miner %+v", miner)
	}
}
