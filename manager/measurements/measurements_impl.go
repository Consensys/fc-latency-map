package measurements

import (
	"time"

	"gorm.io/gorm/clause"

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

func (m *MeasurementServiceImpl) UpsertMeasurements(mrs []*atlas.Measurement) {
	measurements := []*models.Measurement{}
	for _, mr := range mrs {
		mdl := &models.Measurement{
			MeasurementID: mr.ID,
			IsOneOff:      mr.IsOneoff,
			StartTime:     mr.StartTime,
			StopTime:      mr.StopTime,
			Status:        mr.Status.Name,
		}
		if mr.Status.Name == stopped {
			mdl.StatusStopTime = mr.Status.When
		}
		measurements = append(measurements, mdl)
	}

	m.dbCreate(measurements)
}

func (m *MeasurementServiceImpl) GetMeasuresLastResultTime() (measurements []*models.Measurement, measurementsStartTime map[int]int) {
	measurementsStartTime = make(map[int]int)
	measurements = m.GetMeasurements()
	for _, id := range measurements {
		resultTime := m.getLastMeasurementResultTime(id.MeasurementID)
		if resultTime > 0 {
			measurementsStartTime[id.MeasurementID] = resultTime
		}
	}

	return measurements, measurementsStartTime
}

func (m *MeasurementServiceImpl) dbCreate(measurements []*models.Measurement) {
	err := (m.DBMgr).GetDB().
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "measurement_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"status", "status_stop_time", "stop_time"}),
		}).
		Create(&measurements).Error
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
	dbc := (m.DBMgr).GetDB()
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
	dbc = dbc.Model(&models.MeasurementResult{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "probe_id"},
				{Name: "measurement_id"},
				{Name: "measurement_timestamp"},
				{Name: "ip"},
			},
			DoNothing: true,
		}).Create(insert)
	log.WithFields(log.Fields{
		"insert rows": dbc.RowsAffected,
	}).Info("Create measurement MeasurementResults")

	if dbc.Error != nil {
		log.WithFields(log.Fields{
			"err": dbc.Error,
		}).Error("Create measurement MeasurementResults")
	}
}

func (m *MeasurementServiceImpl) GetMeasurements() []*models.Measurement {
	var measurements []*models.Measurement
	dbc := (m.DBMgr).GetDB()
	dbc.Model(&models.Measurement{}).Find(&measurements)

	return measurements
}

func (m *MeasurementServiceImpl) getLastMeasurementResultTime(measurementID int) int {
	dbc := (m.DBMgr).GetDB()

	measurementResults := &models.MeasurementResult{}

	dbc.Model(&models.MeasurementResult{}).
		Select("max(measurement_timestamp) measurement_timestamp").
		Where("measurement_id = ?", measurementID).
		First(&measurementResults)

	return measurementResults.MeasurementTimestamp
}

func (m *MeasurementServiceImpl) GetProbIDs(places []Place, lat, long float64) []string {
	if lat == 0 && long == 0 {
		return []string{}
	}
	p := Place{Latitude: lat, Longitude: long}
	nearestProbesAmount := m.Conf.GetInt("NEAREST_AIRPORTS")
	nearestLocationsIDs := FindNearest(places, p, nearestProbesAmount)
	if len(nearestLocationsIDs) == 0 {
		return []string{}
	}

	var ripeIDs []string
	if err := m.DBMgr.GetDB().
		Select("probes.probe_id").
		Model(&Probes{}).
		Preload("locations").
		Joins("JOIN locations on locations.iata_code=probes.iata_code").
		Where("locations.id in (?)", nearestLocationsIDs).
		Scan(&ripeIDs).Error; err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get probeId from locations")

		return []string{}
	}

	return ripeIDs
}

func (m *MeasurementServiceImpl) PlacesDataSet() ([]Place, error) {
	var places []Place
	err := m.DBMgr.GetDB().Model(&models.Location{}).Find(&places).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get places from db")

		return places, err
	}
	return places, nil
}
