package measurements

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
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

type Probes struct {
	gorm.Model
}

func (m *MeasurementServiceImpl) CreateMeasurements() {
	var miners []*models.Miner
	err := (*m.DbMgr).GetDb().Debug().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Find db miners")
		return
	}

	var probesIDs []string
	// FIXME replace model and probeID column
	err = (*m.DbMgr).GetDb().Debug().Model(models.Miner{}).Select("id").Find(&probesIDs).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Find db miners")
		return
	}

	join := strings.Join(probesIDs, ",")
	_, _ = m.CreatePingProbes(miners, "probes", join)
}
