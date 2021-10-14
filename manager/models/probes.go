package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model      `json:"-"`
	ProbeID         int     `gorm:"column:probe_id;uniqueIndex" json:"probe_id"`
	CountryCode     string  `gorm:"column:country_code" json:"country_code"`
	Status          string  `gorm:"column:status" json:"status"`
	Latitude        float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude       float64 `gorm:"column:longitude" json:"longitude,omitempty"`
	IsAnchor        bool    `gorm:"column:is_anchor" json:"is_anchor"`
	IsPublic        bool    `gorm:"column:is_public" json:"is_public"`
	SwapCoordinates bool    `gorm:"column:swap_coordinates" json:"-"`
}
