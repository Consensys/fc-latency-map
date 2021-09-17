package export

import (
	"time"
)

type Results struct {
	MinersLatency map[string][]*Miners `json:"miners_latency,omitempty"`
}

type Miners struct {
	Address  string        `json:"address"`
	IP       []string      `json:"ip,omitempty"`
	Measures []*MeasuresIP `json:"measures,omitempty"`
}

type MeasuresIP struct {
	IP      string     `json:"ip"`
	Latency []*Latency `json:"latency,omitempty"`
}
type Latency struct {
	Avg  float64   `json:"avg,omitempty"`
	Lts  int       `json:"lts,omitempty"`
	Max  float64   `json:"max,omitempty"`
	Min  float64   `json:"min,omitempty"`
	Date time.Time `json:"date,omitempty"`
	IP   string    `json:"ip,omitempty"`
}
