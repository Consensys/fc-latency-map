package db

import (
	"gorm.io/gorm"
)

type DatabaseMgr interface {

	// GetDB database connection instance
	GetDB() (db *gorm.DB)
}
