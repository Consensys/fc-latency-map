package models

import (
	"time"
)

type ResultsData struct {
	MinersLatency map[string][]*MinersLatency `json:"miners_latency,omitempty"`
}

type MinersLatency struct {
	Address  string        `json:"address"`
	Ip       string        `json:"ip,omitempty"`
	Measures []*MeasuresIp `json:"measures,omitempty"`
}

type MeasuresIp struct {
	Ip           string          `json:"ip"`
	MeasuresData []*MeasuresData `json:"measures_data,omitempty"`
}
type MeasuresData struct {
	Avg  float64   `json:"avg,omitempty"`
	Lts  int       `json:"lts,omitempty"`
	Max  float64   `json:"max,omitempty"`
	Min  float64   `json:"min,omitempty"`
	Date time.Time `json:"date,omitempty"`
	Ip   string    `json:"ip,omitempty"`
}
