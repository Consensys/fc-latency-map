package models

import (
	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	Miner       string
	ProbeID     int
	MeasureDate int
	TimeAverage float64
}
