package measurements

import (
	"time"

	"gorm.io/gorm/clause"

	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type measurementServiceImpl struct {
	Conf  *viper.Viper
	DBMgr db.DatabaseMgr
	FMgr  fmgr.FilecoinMgr
}

func newMeasurementServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr, fMgr fmgr.FilecoinMgr) MeasurementService {
	return &measurementServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
		FMgr:  fMgr,
	}
}

func (m *measurementServiceImpl) UpsertMeasurements(mrs []*atlas.Measurement) {
	var measurements []*models.Measurement
	for _, mr := range mrs {
		mdl := &models.Measurement{
			MeasurementID:  mr.ID,
			IsOneOff:       mr.IsOneoff,
			StartTime:      mr.StartTime,
			StopTime:       mr.StopTime,
			Status:         mr.Status.Name,
			StatusStopTime: mr.Status.When,
		}

		measurements = append(measurements, mdl)
	}

	m.dbCreate(measurements)
}

func (m *measurementServiceImpl) dbCreate(measurements []*models.Measurement) {
	err := (m.DBMgr).GetDB().
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "measurement_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"status", "status_stop_time", "is_one_off", "start_time", "stop_time"}),
		}).
		Create(&measurements).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create Measurement in db")

		return
	}
}

func (m *measurementServiceImpl) GetMinersWithGeolocation() []*models.Miner {
	var miners []*models.Miner

	err := (m.DBMgr).GetDB().
		Where("latitude != 0 and longitude!=0").
		Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMinersWithGeolocation")

		return nil
	}

	return miners
}

func (m *measurementServiceImpl) ImportMeasurement(mr []atlas.MeasurementResult) {
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

func (m *measurementServiceImpl) getMeasurementsRunning() []*models.Measurement {
	var measurements []*models.Measurement
	dbc := (m.DBMgr).GetDB()
	dbc.Model(&models.Measurement{}).
		Find(&measurements, "status in ('', 'running', 'Specified', 'Scheduled', 'Ongoing')")

	return measurements
}

func (m *measurementServiceImpl) getProbIDs(places []Place, lat, long float64) []string {
	if lat == 0 && long == 0 {
		return []string{}
	}
	p := Place{Latitude: lat, Longitude: long}
	nearestProbesAmount := m.Conf.GetInt("NEAREST_AIRPORTS")
	nearestLocationsIDs := FindNearest(places, p, nearestProbesAmount)
	if len(nearestLocationsIDs) == 0 {
		return []string{}
	}

	ripeProbeIDs := []string{}

	dbc := m.DBMgr.GetDB()
	if err := dbc.Model(models.Probe{}).
		Distinct().
		Where("id in (?)",
			dbc.Select("probe_id").
				Table("locations_probes").
				Where("location_id in (?)", nearestLocationsIDs)).
		Pluck("probe_id", &ripeProbeIDs).Error; err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get probeId from locations")

		return []string{}
	}

	return ripeProbeIDs
}

func (m *measurementServiceImpl) getLocationsAsPlaces() ([]Place, error) {
	var places []Place
	err := m.DBMgr.GetDB().
		Model(&models.Location{}).
		Where("latitude!=0 and longitude!=0").
		Find(&places).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get places from db")

		return places, err
	}
	return places, nil
}
