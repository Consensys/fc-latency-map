package measurements

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MeasurementServiceImpl struct {
	Conf  *viper.Viper
	DBMgr db.DatabaseMgr
	FMgr  fmgr.FilecoinMgr
}

func NewMeasurementServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr, fMgr fmgr.FilecoinMgr) MeasurementService {
	return &MeasurementServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
		FMgr:  fMgr,
	}
}

type Probes struct {
	gorm.Model
}

func (m *MeasurementServiceImpl) CreateMeasurements(mrs []*atlas.Measurement) {
	measurements := []*models.Measurement{}
	for _, mr := range mrs {
		measurements = append(measurements,
			&models.Measurement{
				MeasurementID: mr.ID,
				IsOneOff:      mr.IsOneoff,
				StartTime:     mr.StartTime,
				StopTime:      mr.StopTime,
			})
	}

	m.dbCreate(measurements)
}

func (m *MeasurementServiceImpl) GetMeasuresLastResultTime() map[int]int {
	measurements := make(map[int]int)
	for _, id := range m.getRipeMeasurementsID() {
		measurements[id] = m.getLastMeasurementResultTime(id)
	}

	return measurements
}

func (m *MeasurementServiceImpl) dbCreate(measurements []*models.Measurement) {
	err := (m.DBMgr).GetDB().Create(&measurements).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create Measurement in db")

		return
	}
}

func (m *MeasurementServiceImpl) GetMiners() []*models.Miner {
	var miners []*models.Miner

	err := (m.DBMgr).GetDB().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMiners")

		return nil
	}

	return miners
}

func (m *MeasurementServiceImpl) ImportMeasurement(mr []atlas.MeasurementResult) {
	dbc := (m.DBMgr).GetDB().Debug()
	var insert []*models.MeasurementResult
	for _, result := range mr { //nolint:gocritic
		t := time.Unix(int64(result.Timestamp), 0)
		insert = append(insert, &models.MeasurementResult{
			IP:                   result.DstAddr,
			MeasurementID:        result.MsmID,
			ProbeID:              result.PrbID,
			MeasurementTimestamp: result.Timestamp,
			MeasurementDate:      t.Format("2006-01-02"),
			TimeAverage:          result.Avg,
			TimeMax:              result.Max,
			TimeMin:              result.Min,
		})
	}
	affected := dbc.Model(&models.MeasurementResult{}).Create(
		insert).RowsAffected
	log.WithFields(log.Fields{
		"insert rows": affected,
	}).Info("Create measurement MeasurementResults")

	if dbc.Error != nil {
		log.WithFields(log.Fields{
			"err": dbc.Error,
		}).Error("Create measurement MeasurementResults")
	}
}

func (m *MeasurementServiceImpl) getRipeMeasurementsID() []int {
	var ripeIDs []int
	dbc := (m.DBMgr).GetDB().Debug()
	dbc.Model(&models.Measurement{}).Pluck("measurement_id", &ripeIDs)

	return ripeIDs
}

func (m *MeasurementServiceImpl) getLastMeasurementResultTime(measurementID int) int {
	dbc := (m.DBMgr).GetDB().Debug()

	measurementResults := &models.MeasurementResult{}

	dbc.Model(&models.MeasurementResult{}).
		Select("max(measurement_timestamp) measurement_timestamp").
		Where("measurement_id = ?", measurementID).
		First(&measurementResults)

	return measurementResults.MeasurementTimestamp
}

func (m *MeasurementServiceImpl) GetProbIDs(lat, long float64) []string {
	if lat == 0 && long == 0 {
		return []string{}
	}
	p := Place{Latitude: lat, Longitude: long}
	nearestProbesAmount := m.Conf.GetInt("NEAREST_PROBES_AMOUNT")
	nearestProbeIDs := FindNearest(p, nearestProbesAmount, "probes", m.DBMgr.GetDB())

	var ripeIDs []string
	m.DBMgr.GetDB().Debug().Model(&models.Probe{}).
		Select("probe_id").
		Where("id in ?", nearestProbeIDs).
		Find(&ripeIDs)

	return ripeIDs
}
