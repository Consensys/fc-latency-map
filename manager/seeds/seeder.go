package seeds

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	_ "gorm.io/driver/sqlite"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

func Seed() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	err = Execute(dbMgr.GetDb())
	if err != nil {
		log.Fatalf("cannot seed tables: %v", err)
	}
}

// Execute runs the data seed process
func Execute(db *gorm.DB) error {

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2021091516",
			Migrate: func(tx *gorm.DB) error {
				locs := []models.Location{
					{Country: "FR"},
					{Country: "NL"},
					{Country: "PT"},
				}
				return tx.Create(locs).Error
			},
		},
		{
			ID: "2021091517",
			Migrate: func(tx *gorm.DB) error {
				miners := []models.Miner{
					{Address: "f023467", Ip: "151.252.13.181"},
					{Address: "f0694396", Ip: "185.37.217.6,2a04:7340:0:1002::16"},
					{Address: "f022163", Ip: "217.71.253.18,62.171.109.134"},
					{Address: "f01231", Ip: "47.252.15.25,172.17.32.101"},
					{Address: "f01272", Ip: "172.16.117.9"},
					{Address: "f01044351", Ip: "221.144.2.39"},
					{Address: "f0106949", Ip: "192.168.0.200"},
					{Address: "f0149768", Ip: ""},
				}
				return tx.Create(miners).Error
			},
		},
		{
			ID: "2021091518",
			Migrate: func(tx *gorm.DB) error {
				locs := []models.Measurement{
					{
						IsOneoff:      false,
						MeasurementID: 32290390,
						Times:         0,
						StartTime:     1631725577,
						StopTime:      1631725877,
					},
				}
				return tx.Create(locs).Error
			},
		},
		{
			ID: "2021091519",
			Migrate: func(tx *gorm.DB) error {
				locs := []models.Probe{
					{ProbeID: 39, CountryCode: "FR"},
					{ProbeID: 1, CountryCode: "NL"},
					{ProbeID: 24, CountryCode: "PT"},
				}
				return tx.Create(locs).Error
			},
		},
	})

	err := m.Migrate()
	if err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Println("Migrate successfully")

	return nil
}
