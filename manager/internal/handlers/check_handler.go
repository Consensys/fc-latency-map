package handlers

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/geomgr"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

// GetHealthCheckHandler returns system health check
func GetHealthCheckHandler() models.HealthCheck {
	payload := models.HealthCheck{Success: true}
	return payload
}

// GetMetricsHandler returns system metrics
func GetMetricsHandler() models.Metrics {
	// Locations
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	lSer := locations.NewLocationServiceImpl(conf, dbMgr)
	tLocations := lSer.GetTotalLocations()

	// Miners
	nodeURL := conf.GetString("FILECOIN_NODE_URL")
	fMgr, err := fmgr.NewFilecoinImpl(nodeURL)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}

	g := geomgr.NewGeoMgrImpl(conf)
	mSer := miners.NewMinerServiceImpl(conf, dbMgr, fMgr, g)
	tMiners := mSer.GetTotalMiners()

	// Probes
	ripeMgr, err := ripemgr.NewRipeImpl(conf)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	pSer, err := probes.NewProbeServiceImpl(dbMgr, ripeMgr)
	if err != nil {
		panic("failed to start probe service")
	}
	tProbes := pSer.GetTotalProbes()

	payload := models.Metrics{
		Locations: &tLocations,
		Miners:    &tMiners,
		Probes:    &tProbes,
	}
	return payload
}
