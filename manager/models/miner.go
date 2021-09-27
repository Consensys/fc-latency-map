package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Address   string  `gorm:"column:address;uniqueIndex" json:"address"`
	IP        string  `gorm:"column:ip" json:"ip"`
	Latitude  float64 `gorm:"column:latitude" json:"latitude"`
	Longitude float64 `gorm:"column:longitude" json:"longitude"`
}
