package export

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
)

type ExportHandler struct {
	Service *Service
}

func NewExportHandler() *ExportHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	mSer := newExportServiceImpl(conf, dbMgr)

	return &ExportHandler{
		Service: &mSer,
	}
}

func (h *ExportHandler) Export(fn string) {
	(*h.Service).Export(fn)
}
