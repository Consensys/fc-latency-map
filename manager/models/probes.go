package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model
	ProbeID      		int `gorm:"column:probe_id;uniqueIndex"`
	CountryCode     string `gorm:"column:country_code"`
}
