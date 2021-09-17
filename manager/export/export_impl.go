package export

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type ExportServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
}

func NewExportServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr) ExportService {
	return &ExportServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
	}
}

func (m *ExportServiceImpl) export(fn string) {
	measurements := m.GetLatencyMeasurementsStored()

	fullJSON, err := json.MarshalIndent(measurements.MinersLatency, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create json data")
		return
	}

	file.Create(fn, fullJSON)
	log.WithFields(log.Fields{
		"file": fn,
	}).Info("Export successful")
}

func (m *ExportServiceImpl) GetLatencyMeasurementsStored() *Results {
	results := &Results{
		MinersLatency: map[string][]*Miners{},
	}

	loc := m.getLocations()

	for _, l := range loc {
		miners := m.getMiners()

		for _, miner := range miners {
			latency := &Miners{
				Address:  miner.Address,
				Measures: []*MeasuresIP{},
			}
			if miner.IP == "" {
				continue
			}
			latency.IP = strings.Split(miner.IP, ",")
			results.MinersLatency[l.Country] = append(results.MinersLatency[l.Country], latency)
			probes := m.getProbes(l)

			for _, probe := range probes {

				for _, ip := range latency.IP {
					measure := &MeasuresIP{IP: ip}

					meas := m.getMeasureResults(probe, ip)
					if len(meas) > 0 {
						latency.Measures = append(latency.Measures, measure)
					}
					for _, m := range meas {
						measureData := &Latency{
							Avg:  m.TimeAverage,
							Min:  m.TimeMin,
							Max:  m.TimeMax,
							Date: m.MeasurementDate,
						}
						measure.Latency = append(measure.Latency, measureData)
					}
				}
			}
		}
	}

	return results
}

func (m *ExportServiceImpl) getMeasureResults(probe *models.Probe, ip string) []*models.MeasurementResult {
	var meas []*models.MeasurementResult
	err := (*m.DbMgr).GetDb().Debug().Select(
		"ip," +
			"measurement_date," +
			"avg(time_average) time_average," +
			"avg(time_max) time_max," +
			"avg(time_min) time_min").
		Group("ip, measurement_date").
		Where(&models.MeasurementResult{
			ProbeID: probe.ProbeID,
			// MinerAddress: miner.Address,
			IP: ip,
		}).
		Where("measurement_date <  date()").
		Find(&meas).Error

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMeasureResults")
		return nil
	}

	return meas
}

func (m *ExportServiceImpl) getProbes(l *models.Location) []*models.Probe {
	var probes []*models.Probe
	err := (*m.DbMgr).GetDb().Where(&models.Probe{
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

func (m *ExportServiceImpl) getMiners() []*models.Miner {
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

func (m *ExportServiceImpl) getLocations() []*models.Location {
	var loc []*models.Location
	err := (*m.DbMgr).GetDb().
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
