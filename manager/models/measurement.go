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
	ProbeID       int
	MeasurementID int
	MeasureDate   int
	TimeAverage   float64
	TimeMax       float64
	TimeMin       float64
	Ip            string
}
