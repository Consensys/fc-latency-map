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
func Execute(dbc *gorm.DB) error {

	m := gormigrate.New(dbc, gormigrate.DefaultOptions, []*gormigrate.Migration{
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
					{Address: "f023467", IP: "151.252.13.181"},
					{Address: "f0694396", IP: "185.37.217.6,2a04:7340:0:1002::16"},
					{Address: "f022163", IP: "217.71.253.18,62.171.109.134"},
					{Address: "f01231", IP: "47.252.15.25,172.17.32.101"},
					{Address: "f01272", IP: "172.16.117.9"},
					{Address: "f01044351", IP: "221.144.2.39"},
					{Address: "f0106949", IP: "192.168.0.200"},
					{Address: "f0149768", IP: ""},
				}
				return tx.Create(miners).Error
			},
		},
		{
			ID: "2021091518",
			Migrate: func(tx *gorm.DB) error {
				locs := []models.Measurement{
					{MeasurementID: 32290390, StartTime: 1631725577, StopTime: 1631725877},
					{MeasurementID: 32294500, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32294501, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32294502, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32294503, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32294504, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32294505, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32294506, StartTime: 1631785212, StopTime: 1631785512},
					{MeasurementID: 32295718, StartTime: 1631799926, StopTime: 1632404726},
					{MeasurementID: 32295719, StartTime: 1631799926, StopTime: 1632404726},
					{MeasurementID: 32295720, StartTime: 1631799926, StopTime: 1632404726},
					{MeasurementID: 32295721, StartTime: 1631799926, StopTime: 1632404726},
					{MeasurementID: 32295722, StartTime: 1631799926, StopTime: 1632404726},
					{MeasurementID: 32295723, StartTime: 1631799926, StopTime: 1632404726},
					{MeasurementID: 32295724, StartTime: 1631799926, StopTime: 1632404726},
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
