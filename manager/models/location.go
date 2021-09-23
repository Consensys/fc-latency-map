package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Country   string
	IataCode  string `gorm:"index"`
	Latitude  float64
	Longitude float64
}
