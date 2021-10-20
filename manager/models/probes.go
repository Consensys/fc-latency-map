package models

import (
	"gorm.io/gorm"
)

type CoordinatesStatus string
type Status string

const (
	CoordinatesStatusUnknown CoordinatesStatus = "Unknown"
	CoordinatesStatusOk      CoordinatesStatus = "Ok"

	StatusConnected    Status = "Connected"
	StatusDisconnected Status = "Disconnected"
)

type Probe struct {
	gorm.Model        `json:"-"`
	ProbeID           int               `gorm:"column:probe_id;index:idx_probe_probe_id,unique" json:"probe_id"`
	CountryCode       string            `gorm:"column:country_code" json:"country_code"`
	Status            Status            `gorm:"column:status;index:idx_probe_status" json:"status"`
	Latitude          float64           `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude         float64           `gorm:"column:longitude" json:"longitude,omitempty"`
	RipeLatitude      float64           `gorm:"column:ripe_latitude" json:"-"`
	RipeLongitude     float64           `gorm:"column:ripe_longitude" json:"-"`
	IsAnchor          bool              `gorm:"column:is_anchor" json:"is_anchor"`
	IsPublic          bool              `gorm:"column:is_public" json:"is_public"`
	CoordinatesStatus CoordinatesStatus `gorm:"column:coordinates_status;default:Unknown;index:idx_probe_coord_status" json:"coordinates_status"`
	Locations         []*Location       `gorm:"many2many:locations_probes;" json:"-"`
	AddressV4         string
	AddressV6         string
}
