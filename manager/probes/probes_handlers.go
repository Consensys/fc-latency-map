package probes

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
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
	pSer, err := NewProbeServiceImpl(conf, &dbMgr)
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

// Update handle updating probes list
func (pHdl *ProbeHandler) GetAllProbes() []*models.Probe {
	return (*pHdl.PSer).GetAllProbes()
}