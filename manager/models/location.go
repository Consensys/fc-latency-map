package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model `json:"-"`
	Country    string  `json:"country"`
	IataCode   string  `gorm:"index" json:"iata_code"`
	Latitude   float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude  float64 `gorm:"column:longitude" json:"longitude,omitempty"`
}
