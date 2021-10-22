package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model `json:"-"`
	Address    string  `gorm:"column:address;index:idx_address,unique" json:"address"`
	IP         string  `json:"ip,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Port       int     `json:"port,omitempty"`
	Country    string  `json:"-"`
}
