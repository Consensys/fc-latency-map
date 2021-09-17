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
	ripe    *ripemgr.Handler
}

func NewHandler() *Handler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	nodeUrl := conf.GetString("FILECOIN_NODE_URL")
	fMgr, err := fmgr.NewFilecoinImpl(nodeUrl)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}

	mSer := NewMeasurementServiceImpl(conf, &dbMgr, &fMgr)

	return &Handler{
		Service: &mSer,
		ripe:    ripemgr.NewHandler(),
	}
}

func (h *Handler) GetMeasures() {
	measures := (*h.Service).getMeasuresLastResultTime()
	results, err := (*h.ripe).GetMeasurementResults(measures)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("GetMeasurementResults")
		return
	}

	(*h.Service).importMeasurement(results)
}

func (h *Handler) CreateMeasurements() {

	ips := strings.Join((*h.Service).getProbIDs(), ",")
	measures, err := (*h.ripe).CreateMeasurement((*h.Service).getMiners(), ips)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Create Ping")
		return
	}

	(*h.Service).createMeasurements(measures)
}
