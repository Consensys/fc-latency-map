package export

import (
	"encoding/json"
	"strings"

	jg "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type ExportServiceImpl struct {
	Conf  *viper.Viper
	DBMgr db.DatabaseMgr
}

func NewExportServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr) Service {
	return &ExportServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
	}
}

func (m *ExportServiceImpl) export(fn string) {
	measurements := m.GetLatencyMeasurementsStored()

	fullJSON, err := json.MarshalIndent(measurements, "", "  ")
	if err != nil {
		jg.WithFields(jg.Fields{
			"error": err,
		}).Error("Create json data")

		return
	}

	file.Create(fn, fullJSON)
	jg.WithFields(jg.Fields{
		"file": fn,
	}).Info("Export successful")
}

func (m *ExportServiceImpl) GetLatencyMeasurementsStored() *Result {
	results := &Result{
		Measurements: map[string]map[string][]*Miner{},
	}

	loc := m.getLocations()
	miners := m.getMiners()

	for _, l := range loc {
		if _, found := results.Measurements[l.Country]; !found {
			results.Measurements[l.Country] = make(map[string][]*Miner)
		}

		for _, miner := range miners {
			latency := &Miner{
				Address:  miner.Address,
				Measures: []*MeasureIP{},
			}
			if miner.IP == "" {
				continue
			}
			probes := m.getProbes(l)
			latency = m.createLatency(probes, latency, miner.IP)
			results.Measurements[l.Country][l.IataCode] = append(results.Measurements[l.Country][l.IataCode], latency)
		}
	}
	results.Location = loc
	results.Miners = miners
	results.Probes = m.GetAllProbes()

	return results
}

func (m *ExportServiceImpl) GetAllProbes() []*models.Probe {
	var probesList []*models.Probe
	m.DBMgr.GetDB().Find(&probesList)

	return probesList
}

func (m *ExportServiceImpl) createLatency(probes []*models.Probe, latency *Miner, ip string) *Miner {
	for _, probe := range probes {
		for _, ip := range strings.Split(ip, ",") {
			measure := &MeasureIP{IP: ip}

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

	return latency
}

func (m *ExportServiceImpl) getMeasureResults(probe *models.Probe, ip string) []*models.MeasurementResult {
	var meas []*models.MeasurementResult
	err := (m.DBMgr).GetDB().Select(
		"ip," +
			"date(measurement_timestamp, 'unixepoch') measurement_date," +
			"avg(time_average) time_average," +
			"max(time_max) time_max," +
			"min(time_min) time_min").
		Group("ip, measurement_date").
		Where(&models.MeasurementResult{
			ProbeID: probe.ProbeID,
			IP:      ip,
		}).
		Find(&meas).Error
	if err != nil {
		jg.WithFields(jg.Fields{
			"error": err,
		}).Error("GetMeasureResults")

		return nil
	}

	return meas
}

func (m *ExportServiceImpl) getProbes(l *models.Location) []*models.Probe {
	var probes []*models.Probe
	err := (m.DBMgr).GetDB().Where(&models.Probe{
		CountryCode: l.Country,
	}).Find(&probes).Error
	if err != nil {
		jg.WithFields(jg.Fields{
			"error": err,
		}).Error("GetProbes")

		return nil
	}

	return probes
}

func (m *ExportServiceImpl) getMiners() []*models.Miner {
	var miners []*models.Miner

	err := (m.DBMgr).GetDB().Find(&miners).Error
	if err != nil {
		jg.WithFields(jg.Fields{
			"error": err,
		}).Error("GetMiners")

		return nil
	}

	return miners
}

func (m *ExportServiceImpl) getLocations() []*models.Location {
	var loc []*models.Location
	err := (m.DBMgr).GetDB().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}, Desc: false}).
		Find(&loc).Error
	if err != nil {
		jg.WithFields(jg.Fields{
			"error": err,
		}).Error("GetAllLocations")

		return nil
	}

	return loc
}
