package main

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/parser"
)

func main() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	nodeUrl := conf.GetString("FILECOIN_NODE_URL")
	fMgr, err := filecoinmgr.NewFilecoinImpl(nodeUrl)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	mSer := miners.NewMinerServiceImpl(dbMgr, fMgr)
	parser.Parse(conf, fMgr, mSer)
}
