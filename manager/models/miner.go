package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model  `json:"-"`
	Address     string      `gorm:"column:address;uniqueIndex" json:"address"`
	IP          string      `gorm:"column:ip" json:"ip,omitempty"`
	GeoLocation GeoLocation `gorm:"embedded;"`
}

func (m *Miner) BeforeCreate(_ *gorm.DB) (err error) {
	m.GeoLocation.updateTrigonometryLatLong()

	return nil
}

func (m *Miner) BeforeUpdate(_ *gorm.DB) (err error) {
	m.GeoLocation.updateTrigonometryLatLong()

	return nil
}
