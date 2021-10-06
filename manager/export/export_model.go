package export

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type Result struct {
	Locations    []*models.Location             `json:"locations,omitempty"`
	Miners       []*models.Miner                `json:"miners,omitempty"`
	Probes       []*models.Probe                `json:"probes,omitempty"`
	Measurements map[string]map[string][]*Miner `json:"measurements,omitempty"`
}

type Miner struct {
	Address  string       `json:"address"`
	Measures []*MeasureIP `json:"measures,omitempty"`
}

type MeasureIP struct {
	IP      string     `json:"ip"`
	Latency []*Latency `json:"latency,omitempty"`
}

type Latency struct {
	Avg  float64 `json:"avg,omitempty"`
	Lts  int     `json:"lts,omitempty"`
	Max  float64 `json:"max,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Date string  `json:"date,omitempty"`
	IP   string  `json:"ip,omitempty"`
}
