package models

import (
	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Address string
	Ip      string
}
