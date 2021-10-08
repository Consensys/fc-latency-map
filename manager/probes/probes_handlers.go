package probes

import (
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type ProbeHandler struct {
	PSer *ProbeService
}

func NewProbeHandler() *ProbeHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	ripeMgr, err := ripemgr.NewRipeImpl(conf)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	pSer, err := NewProbeServiceImpl(dbMgr, ripeMgr)
	if err != nil {
		panic("failed to start probe service")
	}

	return &ProbeHandler{
		PSer: &pSer,
	}
}

// Update handle updating probes list
func (pHdl *ProbeHandler) Update() {
	(*pHdl.PSer).Update()
}

// List handle updating probes list
func (pHdl *ProbeHandler) List() []*models.Probe {
	return (*pHdl.PSer).ListProbes()
}

// Import handle Import probes from ripe
func (pHdl *ProbeHandler) Import() {
	(*pHdl.PSer).ImportProbes()
}
