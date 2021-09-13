package measurements

import (
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type Handler struct {
	Service *MeasurementService
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

	var apiKey = conf.GetString("RIPE_API_KEY")
	cfg := atlas.Config{
		APIKey: apiKey,
	}
	ripe, err := atlas.NewClient(cfg)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}

	mSer := NewMeasurementServiceImpl(conf, &dbMgr, &fMgr, ripe)

	return &Handler{
		Service: &mSer,
	}
}

func (h *Handler) GetMeasures(s string) {
	(*h.Service).GetRipeMeasures(s)
}

func (h *Handler) CreateMeasurementType(miners []models.Miner, probeType, value string) {
	_, _ = (*h.Service).CreatePingProbes(miners, probeType, value)
}

func (h *Handler) ExportData(fn string) {
	(*h.Service).ExportDbData(fn)
}

func (h *Handler) CreateMeasurements() {
	(*h.Service).CreateMeasurements()
}
