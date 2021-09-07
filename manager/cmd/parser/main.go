package main

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

func main() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	minerService := miners.NewMinerServiceImpl(dbMgr)
	minerService.CreateMiner(&models.Miner{
		Address: "dummyAddress",
		Ip:      "dummyIp",
	})
}
