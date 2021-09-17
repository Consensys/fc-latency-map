package measurements

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/keltia/ripe-atlas"
)

type MeasurementServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
	FMgr  *fmgr.FilecoinMgr
}

func NewMeasurementServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr, fMgr *fmgr.FilecoinMgr) MeasurementService {
	return &MeasurementServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
		FMgr:  fMgr,
	}
}

type Probes struct {
	gorm.Model
}

func (m *MeasurementServiceImpl) createMeasurements(mrs []*atlas.Measurement) {
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

func (m *MeasurementServiceImpl) getMeasuresLastResultTime() map[int]int {
	measurements := make(map[int]int)
	for _, id := range m.getRipeMeasurementsId() {
		measurements[id] = m.getLastMeasurementResultTime(id)
	}
	return measurements
}

func (m *MeasurementServiceImpl) dbCreate(measurements []*models.Measurement) {
	err := (*m.DbMgr).GetDb().Create(&measurements).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create Measurement in db")
		return
	}
}

func (m *MeasurementServiceImpl) getMinersAddress() []string {
	var miners []*models.Miner
	err := (*m.DbMgr).GetDb().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Find miners")
		return nil
	}

	var mAdds = []string{}
	for _, miner := range miners {
		mAdds = append(mAdds, miner.Address)
	}
	return mAdds
}

func (m *MeasurementServiceImpl) getMeasureResults(probe *models.Probe, ip string) []*models.MeasurementResult {
	var meas []*models.MeasurementResult
	err := (*m.DbMgr).GetDb().Debug().Where(&models.MeasurementResult{
		ProbeID: probe.ProbeID,
		Ip:      ip,
	}).Find(&meas).Error

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMeasureResults")
		return nil
	}

	return meas
}

func (m *MeasurementServiceImpl) getProbes(l *models.Location) []*models.Probe {
	var probes []*models.Probe
	err := (*m.DbMgr).GetDb().Debug().Where(&models.Probe{
		CountryCode: l.Country,
	}).Find(&probes).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetProbes")

		return nil
	}
	return probes
}

func (m *MeasurementServiceImpl) getMiners() []*models.Miner {
	var miners []*models.Miner

	err := (*m.DbMgr).GetDb().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMiners")
		return nil
	}
	return miners
}

func (m *MeasurementServiceImpl) importMeasurement(mr []atlas.MeasurementResult) {
	dbc := (*m.DbMgr).GetDb().Debug()

	for _, result := range mr {
		affected := dbc.Model(&models.MeasurementResult{}).Create(&models.MeasurementResult{
			Ip:            result.DstAddr,
			MeasurementID: result.MsmID,
			ProbeID:       result.PrbID,
			MeasureDate:   result.Timestamp,
			TimeAverage:   result.Avg,
			TimeMax:       result.Max,
			TimeMin:       result.Min,
		}).RowsAffected

		log.WithFields(log.Fields{
			"affected": affected,
		}).Info("Create measurement MeasurementResults")

		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"err": dbc.Error,
			}).Error("Create measurement MeasurementResults")
		}
	}
}

func (m *MeasurementServiceImpl) getLocations() []*models.Location {
	var loc []*models.Location
	err := (*m.DbMgr).GetDb().Debug().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}, Desc: false}).
		Find(&loc).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetLocations")

		return nil
	}

	return loc
}

func (m *MeasurementServiceImpl) getRipeMeasurementsId() []int {
	var ripeIDs []int
	dbc := (*m.DbMgr).GetDb().Debug()
	dbc.Model(&models.Measurement{}).Pluck("measurement_id", &ripeIDs)

	return ripeIDs
}

func (m *MeasurementServiceImpl) getLastMeasurementResultTime(measurementID int) int {
	dbc := (*m.DbMgr).GetDb().Debug()

	measurementResults := &models.MeasurementResult{}

	dbc.Model(&models.MeasurementResult{}).
		Select("max(measure_date) measure_date").
		Where("measurement_id = ?", measurementID).
		First(&measurementResults)

	return measurementResults.MeasureDate
}

func (m *MeasurementServiceImpl) getProbIDs() []string {
	var probesIDs []string

	err := (*m.DbMgr).GetDb().Model(models.Probe{}).Select("probe_id").Find(&probesIDs).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Find db miners")

		return nil
	}

	return probesIDs
}
