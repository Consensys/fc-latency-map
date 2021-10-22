package measurements

import (
	"strings"

	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/models"

	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type Handler struct {
	Service MeasurementService
	ripeMgr ripemgr.RipeMgr
}

func NewHandler() *Handler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	nodeURL := conf.GetString("FILECOIN_NODE_URL")
	fMgr, err := fmgr.NewFilecoinImpl(nodeURL)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %stopped", err)
	}

	ripeMgr, err := ripemgr.NewRipeImpl(conf)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %stopped", err)
	}

	mSer := newMeasurementServiceImpl(conf, dbMgr, fMgr)

	return &Handler{
		Service: mSer,
		ripeMgr: ripeMgr,
	}
}
func (h *Handler) ImportMeasures() {
	measurements := h.Service.getMeasurementsRunning()
	for _, m := range measurements {
		measure, err := h.ripeMgr.GetMeasurement(m.MeasurementID)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("GetMeasurement")

			continue
		}
		h.Service.UpsertMeasurements([]*atlas.Measurement{measure})

		log.WithFields(log.Fields{
			"MeasurementID":  m.MeasurementID,
			"StatusStopTime": m.StatusStopTime,
			"StopTime":       m.StopTime,
		}).Info("GetMeasurementResults")

		results, err := h.ripeMgr.GetMeasurementResults(m.MeasurementID)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("GetMeasurementResults")

			continue
		}

		h.Service.ImportMeasurement(results)
	}
}

const oneMinerMeasurement = 1

func (h *Handler) CreateMeasurements(parameters []string) {
	places, err := h.Service.getLocationsAsPlaces()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("placesDataSet")

		return
	}

	miners := h.Service.GetMinersWithGeolocation()
	for i, v := range miners {
		counter := i
		if len(parameters) > 1 {
			if parameters[1] != v.Address {
				continue
			}
			counter = oneMinerMeasurement
		}
		pIDs := strings.Join(h.Service.getProbIDs(places, v.Latitude, v.Longitude), ",")
		log.WithFields(log.Fields{
			"miner.address": v.Address,
			"probeId":       pIDs,
		}).Info("locations Measurements")

		measures, err := h.ripeMgr.CreateMeasurements([]*models.Miner{v}, pIDs, counter)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Create Traceroute")

			continue
		}

		h.Service.UpsertMeasurements(measures)
	}
}
