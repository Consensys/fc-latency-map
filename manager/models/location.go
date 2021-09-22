package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Country   string
	IataCode  string
	Latitude  string
	Longitude string
}
