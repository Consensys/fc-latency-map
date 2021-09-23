package db

//go:generate mockgen -destination mocks.go -package db . DatabaseMgr

import (
	"gorm.io/gorm"
)

type DatabaseMgr interface {

	// GetDB database connection instance
	GetDB() (db *gorm.DB)
}
