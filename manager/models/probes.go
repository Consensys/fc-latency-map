package models

import (
	"gorm.io/gorm"
)

type CoordinatesStatus string
type Status string

const (
	CoordinatesStatusUnknown CoordinatesStatus = "Unknown"
	CoordinatesStatusFixed   CoordinatesStatus = "Fixed"
	CoordinatesStatusOk      CoordinatesStatus = "Ok"

	StatusConnected    Status = "Connected"
	StatusDisconnected Status = "Disconnected"
)

type Probe struct {
	gorm.Model        `json:"-"`
	ProbeID           int               `gorm:"column:probe_id;uniqueIndex" json:"probe_id"`
	CountryCode       string            `gorm:"column:country_code" json:"country_code"`
	Status            Status            `gorm:"column:status;index:idx_probes_status" json:"status"`
	Latitude          float64           `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude         float64           `gorm:"column:longitude" json:"longitude,omitempty"`
	RipeLatitude      float64           `gorm:"column:ripe_latitude" json:"-"`
	RipeLongitude     float64           `gorm:"column:ripe_longitude" json:"-"`
	IsAnchor          bool              `gorm:"column:is_anchor" json:"is_anchor"`
	IsPublic          bool              `gorm:"column:is_public" json:"is_public"`
	CoordinatesStatus CoordinatesStatus `gorm:"column:coordinates_status;index:idx_probes_coord_status;default:Unknown" json:"coordinates_status"`
}
