package measurements

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

func (m *MeasurementServiceImpl) dbCreate(measurements []*models.Measurement) {
	err := (*m.DbMgr).GetDb().Create(&measurements).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create Measurement in db")
		return
	}
}

func (m *MeasurementServiceImpl) dbExportData(fn string) {
	measurements := m.GetLatencyMeasurementsStored()

	fullJson, err := json.MarshalIndent(measurements.MinersLatency, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create json data")
		return
	}

	file.Create(fn, fullJson)
	log.WithFields(log.Fields{
		"file": fn,
	}).Info("Export successful")
}

func (m *MeasurementServiceImpl) GetLatencyMeasurementsStored() *models.ResultsData {
	results := &models.ResultsData{
		MinersLatency: map[string][]*models.MinersLatency{},
	}

	loc := m.getLocations()

	for _, l := range loc {
		miners := m.getMiners()

		for _, miner := range miners {
			latency := &models.MinersLatency{
				Address:  miner.Address,
				Measures: []*models.MeasuresIP{},
			}
			if miner.Ip != "" {
				latency.IP = strings.Split(miner.Ip, ",")
			}
			results.MinersLatency[l.Country] = append(results.MinersLatency[l.Country], latency)
			probes := m.getProbes(l)

			for _, probe := range probes {
				var meas []*models.MeasurementResult
				for _, ip := range strings.Split(miner.Ip, ",") {

					measure := &models.MeasuresIP{IP: ip}

					meas = m.getMeasureResults(probe, miner, ip)
					if len(meas) > 0 {
						latency.Measures = append(latency.Measures, measure)
					}
					for _, m := range meas {
						measureData := &models.MeasuresData{
							Avg:  m.TimeAverage,
							Min:  m.TimeMin,
							Max:  m.TimeMax,
							Date: time.Unix(int64(m.MeasureDate), 0),
						}
						measure.MeasuresData = append(measure.MeasuresData, measureData)
					}
				}
			}
		}
	}

	return results
}

func (m *MeasurementServiceImpl) getMeasureResults(probe *models.Probe, miner *models.Miner, ip string) []*models.MeasurementResult {
	var meas []*models.MeasurementResult
	err := (*m.DbMgr).GetDb().Debug().Where(&models.MeasurementResult{
		ProbeID:      probe.ProbeID,
		MinerAddress: miner.Address,
		Ip:           ip,
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

func (m *MeasurementServiceImpl) importMeasurement(measurementResults []MeasurementResult) {
	db := (*m.DbMgr).GetDb().Debug()
	for _, item := range measurementResults {
		for _, result := range item.Results {
			affected := db.Model(&models.MeasurementResult{}).Create(&models.MeasurementResult{
				Ip:            item.Measurement.Target,
				MeasurementID: item.Measurement.ID,
				MinerAddress:  strings.Join(item.Measurement.Tags, ","),
				ProbeID:       result.PrbID,
				MeasureDate:   result.Timestamp,
				TimeAverage:   result.Avg,
				TimeMax:       result.Max,
				TimeMin:       result.Min,
			}).RowsAffected

			log.WithFields(log.Fields{
				"affected": affected,
			}).Info("Create measurement MeasurementResults")

			if db.Error != nil {
				log.WithFields(log.Fields{
					"err": db.Error,
				}).Error("Create measurement MeasurementResults")
			}
		}
	}
}

func (m *MeasurementServiceImpl) getRipeMeasurementsId() []int {
	db := (*m.DbMgr).GetDb().Debug()
	var ripeIDs []int
	db.Model(&models.Measurement{}).Pluck("measurement_id", &ripeIDs)
	return ripeIDs
}

func (m *MeasurementServiceImpl) getLastMeasurementResultTime(measurementID int) int {
	db := (*m.DbMgr).GetDb().Debug()

	measurementResults := &models.MeasurementResult{}
	db.Model(&models.MeasurementResult{}).
		Select("max(measure_date) measure_date").
		Where("measurement_id = ?", measurementID).
		First(&measurementResults)

	return measurementResults.MeasureDate
}
