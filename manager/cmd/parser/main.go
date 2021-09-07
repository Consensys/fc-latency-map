package main

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/model"
)

func main() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	db := dbMgr.GetDb()
	db.Create(&model.Miner{
		Address: "dummyAddress",
		Ip:      "dummyIp",
	})
}
