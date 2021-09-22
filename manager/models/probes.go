package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model
	ProbeID     int     `gorm:"column:probe_id;uniqueIndex"`
	CountryCode string  `gorm:"column:country_code"`
	Latitude    float64 `gorm:"column:latitude"`
	Longitude   float64 `gorm:"column:longitude"`
}
