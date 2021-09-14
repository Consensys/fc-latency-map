package models

import (
	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	ProbeID      int
	MeasureDate  int
	TimeAverage  float64
	TimeMax      float64
	TimeMin      float64
	MinerAddress string
	Ip           string
}
