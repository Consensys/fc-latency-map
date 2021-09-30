package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model  `json:"-"`
	ProbeID     int     `gorm:"column:probe_id;uniqueIndex" json:"probe_id"`
	CountryCode string  `gorm:"column:country_code" json:"country_code"`
	Latitude    float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude   float64 `gorm:"column:longitude" json:"longitude,omitempty"`
}
