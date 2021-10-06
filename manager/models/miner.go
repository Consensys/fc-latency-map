package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model `json:"-"`
	Address    string  `gorm:"column:address;uniqueIndex" json:"address"`
	IP         string  `json:"ip,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}
