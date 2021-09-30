package models

import (
	"math"

	"github.com/ConsenSys/fc-latency-map/manager/radians"
)

type GeoLocation struct {
	Latitude     float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude    float64 `gorm:"column:longitude" json:"longitude,omitempty"`
	CosLatitude  float64 ` json:"-"`
	SinLatitude  float64 ` json:"-"`
	CosLongitude float64 ` json:"-"`
	SinLongitude float64 ` json:"-"`
}

func (m *GeoLocation) updateTrigonometryLatLong() {
	m.SinLatitude = math.Sin(radians.Radians(m.Latitude))
	m.CosLatitude = math.Cos(radians.Radians(m.Latitude))

	m.SinLongitude = math.Sin(radians.Radians(m.Longitude))
	m.CosLongitude = math.Cos(radians.Radians(m.Longitude))
}
