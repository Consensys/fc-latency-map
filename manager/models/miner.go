package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Address   string  `gorm:"column:address;uniqueIndex"`
	IP        string  `gorm:"column:ip"`
	Latitude  float64 `gorm:"column:latitude"`
	Longitude float64 `gorm:"column:longitude"`
}
