package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type DatabaseMgrMock struct {
	db *gorm.DB
}

func NewMockDatabaseMgr() DatabaseMgr {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	runMigrate(db)
	return &DatabaseMgrMock{
		db: db,
	}
}

func runMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(&models.Miner{})
	_ = db.AutoMigrate(&models.Location{})
	_ = db.AutoMigrate(&models.Measurement{})
	_ = db.AutoMigrate(&models.MeasurementResult{})
	_ = db.AutoMigrate(&models.Probe{})
}

func (dbMgr *DatabaseMgrMock) GetDB() *gorm.DB {
	return dbMgr.db
}
