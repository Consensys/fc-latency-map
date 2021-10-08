package measurements

import (
	"strconv"
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
	locations := make(map[int][]*models.Location)

	for _, result := range mr { //nolint:gocritic
		t := time.Unix(int64(result.Timestamp), 0)

		if _, found := locations[result.PrbID]; !found {
			locations[result.PrbID] = m.getLocationsWithProbID(result.PrbID)
		}

		insert = append(insert, &models.MeasurementResult{
			IP:                   result.DstAddr,
			MeasurementID:        result.MsmID,
			ProbeID:              result.PrbID,
			MeasurementTimestamp: result.Timestamp,
			MeasurementDate:      t.Format("2006-01-02"),
			TimeAverage:          result.Avg,
			TimeMax:              result.Max,
			TimeMin:              result.Min,
			Locations:            locations[result.PrbID],
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

func (m *measurementServiceImpl) GetMeasurementsRunning() []*models.Measurement {
	var measurements []*models.Measurement
	dbc := (m.DBMgr).GetDB()
	dbc.Model(&models.Measurement{}).
		Find(&measurements, "status not in ('Failed', 'Stopped')")

	return measurements
}

func (m *measurementServiceImpl) GetProbIDs(places []Place, lat, long float64) []string {
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
	var nearestLocations []*models.Location

	if err := m.DBMgr.GetDB().Model(models.Location{}).
		Preload(clause.Associations).
		Find(&nearestLocations, "id in ?", nearestLocationsIDs).
		Error; err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get probeId from locations")

		return []string{}
	}
	for _, location := range nearestLocations {
		for _, probe := range location.Probes {
			ripeProbeIDs = add(ripeProbeIDs, strconv.Itoa(probe.ProbeID))
		}
	}

	return ripeProbeIDs
}

func add(s []string, str string) []string {
	for _, v := range s {
		if v == str {
			return s
		}
	}

	return append(s, str)
}

func (m *measurementServiceImpl) GetLocationsAsPlaces() ([]Place, error) {
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

func (m *measurementServiceImpl) getLocationsWithProbID(ripeProbeID int) []*models.Location {
	l := []*models.Location{}
	dbc := m.DBMgr.GetDB()
	err := dbc.Where("id in (?)", dbc.
		Select("location_id").
		Table("locations_probes").
		Where("probe_id in (?)",
			dbc.Select("id").
				Table("probes").
				Where("probe_id in (?)", ripeProbeID)),
	).Find(&l).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("getLocationsWithProbID")
	}
	return l
}
