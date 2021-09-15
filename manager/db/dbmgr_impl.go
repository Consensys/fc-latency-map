package db

import (
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type DatabaseMgrImpl struct {
	Db *gorm.DB
}

func NewDatabaseMgrImpl(conf *viper.Viper) (DatabaseMgr, error) {
	dbName := conf.GetString("DB_CONNECTION")
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	migrate(db)
	return &DatabaseMgrImpl{
		Db: db,
	}, nil
}

func migrate(db *gorm.DB) {
	_ = db.AutoMigrate(&models.Miner{})
	_ = db.AutoMigrate(&models.Location{})
	_ = db.AutoMigrate(&models.Measurement{})
	_ = db.AutoMigrate(&models.MeasurementResult{})
	_ = db.AutoMigrate(&models.Probe{})
}

func (dbMgr *DatabaseMgrImpl) GetDb() *gorm.DB {
	return dbMgr.Db
}
