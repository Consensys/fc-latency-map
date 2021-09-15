package models

import (
	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	IsOneoff      bool
	MeasurementID int
	Times         int
	StartTime     int
	StopTime      int
}

type MeasurementResults struct {
	gorm.Model
	ProbeID       int
	MeasurementID int
	MeasureDate   int
	TimeAverage   float64
	TimeMax       float64
	TimeMin       float64
	MinerAddress  string
	Ip            string
}
