package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Address string `gorm:"column:address;unique"`
	Ip      string `gorm:"column:ip"`
}
