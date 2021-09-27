package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Country   string  `json:"country"`
	IataCode  string  `gorm:"index" json:"iata_code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
