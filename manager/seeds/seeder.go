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
					{MeasurementID: 32523262},
					{MeasurementID: 32517958},
				}
				return tx.Create(locs).Error
			},
		},
		{
			ID: "2021091519",
			Migrate: func(tx *gorm.DB) error {
				locs := []models.Probe{
					{ProbeID: 1001079},
					{ProbeID: 1000555},
					{ProbeID: 1000916},
					{ProbeID: 1001066},
					{ProbeID: 1000555},
					{ProbeID: 1001268},
					{ProbeID: 10019},
					{ProbeID: 1002886},
					{ProbeID: 1002931},
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
