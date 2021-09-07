package db

import (
	"gorm.io/gorm"
)

type DatabaseMgr interface {

	// Insert a new value
	Create(value interface{}) (tx *gorm.DB)
}
