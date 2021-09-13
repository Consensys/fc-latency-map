package measurements

import (
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/keltia/ripe-atlas"
)

type MeasurementResult struct {
	Measurement atlas.Measurement
	Results     []atlas.MeasurementResult
}

type MeasurementServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
	FMgr  *fmgr.FilecoinMgr
	Ripe  *atlas.Client
}

func NewMeasurementServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr, fMgr *fmgr.FilecoinMgr, r *atlas.Client) MeasurementService {
	return &MeasurementServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
		FMgr:  fMgr,
		Ripe:  r,
	}
}
