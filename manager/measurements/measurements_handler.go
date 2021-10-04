package measurements

import (
	"strings"

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
		log.Fatalf("connecting with lotus failed: %s", err)
	}

	ripeMgr, err := ripemgr.NewRipeImpl(conf)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}

	mSer := NewMeasurementServiceImpl(conf, dbMgr, fMgr)

	return &Handler{
		Service: mSer,
		ripeMgr: ripeMgr,
	}
}

func (h *Handler) ImportMeasures() {
	measures := h.Service.GetMeasuresLastResultTime()
	results, err := h.ripeMgr.GetMeasurementResults(measures)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("GetMeasurementResults")

		return
	}

	h.Service.ImportMeasurement(results)
}

func (h *Handler) CreateMeasurements() {
	places, err := h.Service.PlacesDataSet()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("placesDataSet")

		return
	}

	miners := h.Service.GetMiners()
	for i, v := range miners {
		pIDs := strings.Join(h.Service.GetProbIDs(places, v.Latitude, v.Longitude), ",")
		log.WithFields(log.Fields{
			"miner.address": v.Address,
			"probeId":       pIDs,
		}).Info("locations Measurements")

		measures, err := h.ripeMgr.CreateMeasurements([]*models.Miner{v}, pIDs, i)

		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Create Ping")

			continue
		}

		h.Service.CreateMeasurements(measures)
	}
}
