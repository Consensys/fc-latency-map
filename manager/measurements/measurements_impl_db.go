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

func (m *MeasurementServiceImpl) GetLatencyMeasurementsStored() *models.ResultsData {

	results := &models.ResultsData{
		MinersLatency: map[string][]*models.MinersLatency{},
	}
	var loc []*models.Location

	(*m.DbMgr).GetDb().Debug().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}, Desc: false}).
		Find(&loc)

	for _, l := range loc {

		var miners []*models.Miner
		err := (*m.DbMgr).GetDb().Find(&miners).Error

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Find miners")
			return nil
		}

		for _, miner := range miners {
			latency := &models.MinersLatency{
				Address:  miner.Address,
				Ip:       strings.Split(miner.Ip, ","),
				Measures: []*models.MeasuresIp{},
			}
			results.MinersLatency[l.Country] = append(results.MinersLatency[l.Country], latency)
			var probes []*models.Probe
			(*m.DbMgr).GetDb().Debug().Where(&models.Probe{
				CountryCode: l.Country,
			}).Find(&probes)

			for _, probe := range probes {
				var meas []*models.MeasurementResults
				for _, ip := range strings.Split(miner.Ip, ",") {

					measure := &models.MeasuresIp{
						Ip: ip,
					}

					(*m.DbMgr).GetDb().Debug().Where(&models.MeasurementResults{
						ProbeID:      probe.ProbeID,
						MinerAddress: miner.Address,
						Ip:           ip,
					}).Find(&meas)
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

func (m *MeasurementServiceImpl) importMeasurement(measurementResults []MeasurementResult) {
	db := (*m.DbMgr).GetDb().Debug()
	for _, item := range measurementResults {
		for _, result := range item.Results {
			affected := db.Model(&models.MeasurementResults{}).Create(&models.MeasurementResults{
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

	measurementResults := &models.MeasurementResults{}
	db.Model(&models.MeasurementResults{}).
		Select("max(measure_date) measure_date").
		Where("measurement_id = ?", measurementID).
		First(&measurementResults)

	return measurementResults.MeasureDate
}
