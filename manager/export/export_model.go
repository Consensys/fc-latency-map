package export

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type Result struct {
	Locations    []*models.Location             `json:"locations,omitempty"`
	Miners       []*models.Miner                `json:"miners,omitempty"`
	Dates        []string                       `json:"dates,omitempty"`
	Measurements map[string]map[string][]*Miner `json:"measurements,omitempty"`
}

type Miner struct {
	Address  string       `json:"address"`
	Measures []*MeasureIP `json:"measures,omitempty"`
}

type MeasureIP struct {
	IP              string  `json:"ip"`
	Avg             float64 `gorm:"column:time_average" json:"avg"`
	MeasurementDate string  `json:"-"`
}
