package db

import (
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
	return &DatabaseMgrImpl{
		Db: db,
	}, nil
}

func (dbMgr *DatabaseMgrImpl) GetDb() (db *gorm.DB) {
	return dbMgr.Db
}
