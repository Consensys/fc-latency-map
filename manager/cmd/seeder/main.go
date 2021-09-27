package main

import (
	log "github.com/sirupsen/logrus"

	_ "gorm.io/driver/sqlite"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/seeds"
)

func main() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	err = seeds.Execute(dbMgr.GetDB())
	if err != nil {
		log.Fatalf("cannot seed tables: %v", err)
	}

	var locations []models.Location
	var count int64
	dbMgr.GetDB().Model(&locations).Count(&count)
	log.Printf("Total locations: %d\n", count)
}
