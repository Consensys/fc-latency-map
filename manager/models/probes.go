package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model
	MinerID      		int `gorm:"column:miner_id;uniqueIndex"`
	CountryCode     string `gorm:"column:country_code"`
}
