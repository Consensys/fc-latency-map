package models

import (
	"gorm.io/gorm"
)

type Probe struct {
	gorm.Model  `json:"-"`
	ProbeID     int         `gorm:"column:probe_id;uniqueIndex" json:"probe_id"`
	CountryCode string      `gorm:"column:country_code" json:"country_code"`
	GeoLocation GeoLocation `gorm:"embedded;"`
}

func (m *Probe) BeforeCreate(_ *gorm.DB) (err error) {
	m.GeoLocation.updateTrigonometryLatLong()

	return nil
}

func (m *Probe) BeforeUpdate(_ *gorm.DB) (err error) {
	m.GeoLocation.updateTrigonometryLatLong()

	return nil
}
