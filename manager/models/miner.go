package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model `json:"-"`
	Address    string  `gorm:"column:address;index:idx_miner_address,unique" json:"address"`
	IP         string  `gorm:"column:ip;index:idx_miner_ip" json:"ip,omitempty"`
	Latitude   float64 `gorm:"column:latitude;index:idx_miner_lat" json:"latitude,omitempty"`
	Longitude  float64 `gorm:"column:longitude;index:idx_miner_long" json:"longitude,omitempty"`
	Port       int     `json:"port,omitempty"`
	Country    string  `json:"-"`
}
