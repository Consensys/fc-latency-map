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

func (m *MeasurementServiceImpl) RipeCreateMeasurements() {
	var miners []*models.Miner
	err := (*m.DbMgr).GetDb().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Find db miners")
		return
	}

	var probesIDs []string
	err = (*m.DbMgr).GetDb().Model(models.Probe{}).Select("probe_id").Find(&probesIDs).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Find db miners")
		return
	}

	join := strings.Join(probesIDs, ",")
	mr, p, err := m.RipeCreatePingWithProbes(miners, join)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Create Ping")
		return
	}

	measurements := []*models.Measurement{}
	for i := range mr.Definitions {
		measurements = append(measurements,
			&models.Measurement{
				MeasurementID: p.Measurements[i],
				IsOneoff:      mr.IsOneoff,
				Times:         mr.Times,
				StartTime:     mr.StartTime,
				StopTime:      mr.StopTime,
			})
	}

	m.dbCreate(measurements)
}

func (m *MeasurementServiceImpl) RipeGetMeasures() {

	for _, id := range m.getRipeMeasurementsId() {

		start := m.getLastMeasurementResultTime(id)

		measurementResults, err := m.getRipeMeasurementResultsById(id, start)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("Load measurement MeasurementResults from Ripe")
		}

		m.importMeasurement(measurementResults)

		log.Info("measurements successfully get")
	}

}
