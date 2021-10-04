package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model  `json:"-"`
	ProbeID     int      `gorm:"column:probe_id;uniqueIndex" json:"probe_id"`
	CountryCode string   `gorm:"column:country_code" json:"country_code"`
	IataCode    string   `gorm:"foreignKey:iata_code" json:"iata_code"`
	Location    Location `gorm:"foreignkey:IataCode;references:iata_code" json:"-"`
	Latitude    float64  `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude   float64  `gorm:"column:longitude" json:"longitude,omitempty"`
}
