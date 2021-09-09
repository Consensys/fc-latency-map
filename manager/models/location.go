package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Country 	string
	Latitude  string
	Longitude string
}
