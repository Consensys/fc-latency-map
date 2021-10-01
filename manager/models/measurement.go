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
	ProbeID              int         `gorm:"foreignKey:probe_id"`
	Probe                Probe       `gorm:"foreignkey:ProbeID;references:probe_id"`
	MeasurementID        int         `gorm:"foreignKey:measurement_id;index"`
	Measurement          Measurement `gorm:"foreignkey:MeasurementID;references:measurement_id"`
	MeasurementTimestamp int         `gorm:"index"`
	IP                   string      `gorm:"index"`
	MeasurementDate      string
	TimeAverage          float64
	TimeMax              float64
	TimeMin              float64
}
