package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model `json:"-"`
	Address    string  `gorm:"column:address;uniqueIndex" json:"address"`
	IP         string  `gorm:"column:ip" json:"ip,omitempty"`
	Latitude   float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude  float64 `gorm:"column:longitude" json:"longitude,omitempty"`
}
