package export

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type ServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
}

func NewServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr) Service {
	return &ServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
	}
}

func (m *ServiceImpl) export(fn string) {
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

func (m *ServiceImpl) GetLatencyMeasurementsStored() *Results {
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
			if miner.Ip == "" {
				continue
			}
			latency.IP = strings.Split(miner.Ip, ",")
			results.MinersLatency[l.Country] = append(results.MinersLatency[l.Country], latency)
			probes := m.getProbes(l)

			for _, probe := range probes {
				var meas []*models.MeasurementResult
				for _, ip := range latency.IP {

					measure := &MeasuresIP{IP: ip}

					meas = m.getMeasureResults(probe, ip)
					if len(meas) > 0 {
						latency.Measures = append(latency.Measures, measure)
					}
					for _, m := range meas {
						measureData := &Latency{
							Avg:  m.TimeAverage,
							Min:  m.TimeMin,
							Max:  m.TimeMax,
							Date: time.Unix(int64(m.MeasureDate), 0),
						}
						measure.Latency = append(measure.Latency, measureData)
					}
				}
			}
		}
	}

	return results
}

func (m *ServiceImpl) getMeasureResults(probe *models.Probe, ip string) []*models.MeasurementResult {
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

func (m *ServiceImpl) getProbes(l *models.Location) []*models.Probe {
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

func (m *ServiceImpl) getMiners() []*models.Miner {
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

func (m *ServiceImpl) getLocations() []*models.Location {
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
