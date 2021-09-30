package seeds

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

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

	err = Execute(dbMgr.GetDB())
	if err != nil {
		log.Fatalf("cannot seed tables: %v", err)
	}
}

// nolint
// Execute runs the data seed process
func Execute(dbc *gorm.DB) error {
	m := gormigrate.New(dbc, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2021091516",
			Migrate: func(tx *gorm.DB) error {
				locs := []models.Location{
					{Country: "FR", IataCode: "CDG", Latitude: 49.012798, Longitude: 2.55},
					{Country: "NL", IataCode: "AMS", Latitude: 52.308601, Longitude: 4.76389},
					{Country: "PT", IataCode: "LIS", Latitude: 38.7813, Longitude: -9.13592},
					{Country: "CN", IataCode: "PEK", Latitude: 40.080101013183594, Longitude: 116.58499908447266},
				}
				return tx.Create(locs).Error
			},
		},
		{
			ID: "2021091517",
			Migrate: func(tx *gorm.DB) error {
				miners := []models.Miner{
					{Address: "f0694396", IP: "185.37.217.6,2a04:7340:0:1002::16", Latitude: 52.48395, Longitude: -1.88980},
					{Address: "f022163", IP: "217.71.253.18,62.171.109.134", Latitude: 47.36329, Longitude: 8.55014},
					{Address: "f01231", IP: "47.252.15.25,172.17.32.101", Latitude: 37.55983, Longitude: -122.27148},
					{Address: "f01044351", IP: "221.144.2.39", Latitude: 37.41043, Longitude: 127.13716},
					{Address: "f01272", IP: "172.16.117.9"},
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
					{ProbeID: 39, CountryCode: "FR", Latitude: 43.2915, Longitude: 1.6185},
					{ProbeID: 1, CountryCode: "NL", Latitude: 52.3475, Longitude: 4.9275},
					{ProbeID: 24, CountryCode: "PT", Latitude: 38.7295, Longitude: -9.1515},
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
