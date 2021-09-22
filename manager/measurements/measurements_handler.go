package measurements

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type Handler struct {
	Service *MeasurementService
	ripeMgr *ripemgr.RipeMgr
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
		Service: &mSer,
		ripeMgr: &ripeMgr,
	}
}

func (h *Handler) GetMeasures() { //nolint:revive
	measures := (*h.Service).getMeasuresLastResultTime()
	results, err := (*h.ripeMgr).GetMeasurementResults(measures)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("GetMeasurementResults")
		return
	}

	(*h.Service).importMeasurement(results)
}

func (h *Handler) CreateMeasurements() {
	pIDs := strings.Join((*h.Service).getProbIDs(), ",")
	measures, err := (*h.ripeMgr).CreateMeasurements((*h.Service).getMiners(), pIDs)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Create Ping")
		return
	}

	(*h.Service).createMeasurements(measures)
}
