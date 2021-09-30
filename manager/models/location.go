package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model  `json:"-"`
	Country     string      `json:"country"`
	IataCode    string      `gorm:"index" json:"iata_code"`
	GeoLocation GeoLocation `gorm:"embedded;"`
}

func (m *Location) BeforeCreate(_ *gorm.DB) (err error) {
	m.GeoLocation.updateTrigonometryLatLong()

	return nil
}

func (m *Location) BeforeUpdate(_ *gorm.DB) (err error) {
	m.GeoLocation.updateTrigonometryLatLong()

	return nil
}
