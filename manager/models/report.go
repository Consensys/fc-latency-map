package models

import (
	"time"
)

type Latency struct {
	Address   string          `json:"address"`
	Locations []*LocationData `json:"locations,omitempty"`
	Ip        string          `json:"ip,omitempty"`
}

type LocationData struct {
	Country   string         `json:"country,omitempty"`
	Latitude  string         `json:"latitude,omitempty"`
	Longitude string         `json:"longitude,omitempty"`
	Measures  []*MeasureData `json:"measures,omitempty"`
}

type MeasureData struct {
	Avg  float64   `json:"avg,omitempty"`
	Lts  int       `json:"lts,omitempty"`
	Max  float64   `json:"max,omitempty"`
	Min  float64   `json:"min,omitempty"`
	Date time.Time `json:"date,omitempty"`
	Ip   string    `json:"ip,omitempty"`
}
