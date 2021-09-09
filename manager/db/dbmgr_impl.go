package db

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	db.AutoMigrate(&models.Miner{})
	db.AutoMigrate(&models.Location{})
}

func (dbMgr *DatabaseMgrImpl) GetDb() (db *gorm.DB) {
	return dbMgr.Db
}
