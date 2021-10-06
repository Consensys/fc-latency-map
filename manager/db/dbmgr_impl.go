package db

import (
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type DatabaseMgrImpl struct {
	db *gorm.DB
}

func NewDatabaseMgrImpl(conf *viper.Viper) (DatabaseMgr, error) {
	dbName := conf.GetString("DB_CONNECTION")
	const sqliteMaxVariables = 999
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		CreateBatchSize: sqliteMaxVariables,
		Logger:          logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	if conf.GetBool("SQL_DEBUG") {
		db = db.Debug()
	}
	migrate(db)

	return &DatabaseMgrImpl{
		db: db,
	}, nil
}

func migrate(db *gorm.DB) {
	_ = db.AutoMigrate(&models.Miner{})
	_ = db.AutoMigrate(&models.Location{})
	_ = db.AutoMigrate(&models.Measurement{})
	_ = db.AutoMigrate(&models.MeasurementResult{})
	_ = db.AutoMigrate(&models.Probe{})
}

func (dbMgr *DatabaseMgrImpl) GetDB() *gorm.DB {
	return dbMgr.db
}
