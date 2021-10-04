package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model `json:"-"`
	Name       string
	Country    string
	IataCode   string  `gorm:"uniqueIndex"`
	Latitude   float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude  float64 `gorm:"column:longitude" json:"longitude,omitempty"`
	Type       string
}
