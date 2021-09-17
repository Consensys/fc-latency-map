package export

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
)

type Handler struct {
	Service *Service
}

func NewHandler() *Handler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	mSer := NewServiceImpl(conf, &dbMgr)

	return &Handler{
		Service: &mSer,
	}
}

func (h *Handler) ExportData(fn string) {
	(*h.Service).export(fn)
}
