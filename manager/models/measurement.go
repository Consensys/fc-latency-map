package models

import (
	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	IsOneOff      bool
	MeasurementID int
	StartTime     int
	StopTime      int
}

type MeasurementResult struct {
	gorm.Model
	ProbeID              int
	MeasurementID        int `gorm:"index"`
	MeasurementTimestamp int `gorm:"index"`
	MeasurementDate      string
	TimeAverage          float64
	TimeMax              float64
	TimeMin              float64
	IP                   string `gorm:"index"`
}
