package db

import (
	"gorm.io/gorm"
)

type DatabaseMgr interface {

	// Get database connection instance
	GetDb() (db *gorm.DB)
}
