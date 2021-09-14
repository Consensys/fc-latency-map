package measurements

import (
	"encoding/json"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

func (m *MeasurementServiceImpl) ExportDbData(fn string) {
	measurements := m.GetLatencyMeasurementsStored()

	fullJson, err := json.MarshalIndent(measurements, "", "  ")
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

func (m *MeasurementServiceImpl) GetLatencyMeasurementsStored() []*models.Latency {
	var miners []*models.Miner
	err := (*m.DbMgr).GetDb().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Find miners")
		return nil
	}

	var latencies = []*models.Latency{}
	for _, miner := range miners {
		latency := &models.Latency{
			Address: miner.Address,
			Ip:      strings.Split(miner.Ip, ","),
		}
		latencies = append(latencies, latency)

		var loc []*models.Location
		(*m.DbMgr).GetDb().Debug().Find(&loc)

		for _, l := range loc {
			location := &models.LocationData{
				Country:   l.Country,
				Longitude: l.Longitude,
				Latitude:  l.Latitude,
			}
			latency.Locations = append(latency.Locations, location)
			for meas := 0; meas < 5; meas++ {
				measureData := &models.MeasureData{
					Avg: rand.Float64(),
					Lts: rand.Int(),
					Max: rand.Float64(),
					Min: rand.Float64(),
				}
				location.Measures = append(location.Measures, measureData)
			}
		}

	}

	return latencies
}

func (m *MeasurementServiceImpl) importMeasurement(measurementResults []MeasurementResult) {

	db := (*m.DbMgr).GetDb().Debug()
	for _, item := range measurementResults {
		for _, result := range item.Results {
			affected := db.Model(&models.Measurement{}).Create(&models.Measurement{
				Miner:       item.Measurement.Target,
				ProbeID:     result.PrbID,
				MeasureDate: result.Timestamp,
				TimeAverage: result.Avg,
			}).RowsAffected

			log.WithFields(log.Fields{
				"affected": affected,
			}).Info("Create measurement Results")

			if db.Error != nil {
				log.WithFields(log.Fields{
					"err": db.Error,
				}).Error("Create measurement Results")
			}
		}
	}
}
