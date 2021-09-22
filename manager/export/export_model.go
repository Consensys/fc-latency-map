package export

type Result struct {
	MinersLatency map[string][]*Miner `json:"minersLatency,omitempty"`
}

type Miner struct {
	Address  string       `json:"address"`
	IP       []string     `json:"ip,omitempty"`
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
