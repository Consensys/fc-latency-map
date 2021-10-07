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
	DBMgr db.DatabaseMgr
}

func newExportServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr) Service {
	return &ExportServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
	}
}

func (m *ExportServiceImpl) export(fn string) {
	measurements := m.getLatencyMeasurementsStored()

	fullJSON, err := json.MarshalIndent(measurements, "", "  ")
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

func (m *ExportServiceImpl) getLatencyMeasurementsStored() *Result {
	var iataCodes []string

	results := &Result{Measurements: map[string]map[string][]*Miner{}}
	loc := m.getLocations()
	miners := m.getMiners()

	for _, l := range loc {
		for _, miner := range miners {
			latency := &Miner{
				Address:  miner.Address,
				Measures: []*MeasureIP{},
			}
			if miner.IP == "" {
				continue
			}

			latency = m.getLatency(l.Probes, latency, miner.IP)
			if len(latency.Measures) > 0 {
				iataCodes = addNewString(iataCodes, l.IataCode)
				if _, found := results.Measurements[l.Country]; !found {
					results.Measurements[l.Country] = make(map[string][]*Miner)
				}
				results.Measurements[l.Country][l.IataCode] = append(results.Measurements[l.Country][l.IataCode], latency)
			}
		}
	}
	m.addRootData(results, miners, iataCodes)

	return results
}

func (m *ExportServiceImpl) addRootData(results *Result, miners []*models.Miner, iataCodes []string) {
	results.Miners = miners
	results.Dates = m.getDates()
	results.Locations = m.getLocationsFromIata(iataCodes)
	for _, location := range results.Locations {
		results.Probes = appendNewProbe(results.Probes, location.Probes)
	}
}

func (m *ExportServiceImpl) getLatency(probes []*models.Probe, latency *Miner, ip string) *Miner {
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
	err := m.DBMgr.GetDB().Select(
		"ip," +
			"measurement_date," +
			"avg(time_average) time_average," +
			"max(time_max) time_max," +
			"min(time_min) time_min").
		Group("ip, measurement_date").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "measurement_date"}, Desc: true}).
		Where(&models.MeasurementResult{
			ProbeID: probe.ProbeID,
			IP:      ip,
		}).
		Find(&meas).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMeasureResults")

		return nil
	}

	return meas
}

func (m *ExportServiceImpl) getMiners() []*models.Miner {
	var miners []*models.Miner

	err := m.DBMgr.GetDB().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetMinersWithGeoLocation")

		return nil
	}

	return miners
}

func (m *ExportServiceImpl) getLocations() []*models.Location {
	var loc []*models.Location
	err := m.DBMgr.GetDB().
		Preload(clause.Associations).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}, Desc: false}).
		Find(&loc).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetAllLocations")

		return nil
	}

	return loc
}

func (m *ExportServiceImpl) getLocationsFromIata(codes []string) []*models.Location {
	var loc []*models.Location
	err := m.DBMgr.GetDB().
		Preload(clause.Associations).
		Where("iata_code in ?", codes).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}, Desc: false}).
		Find(&loc).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetAllLocations")

		return nil
	}

	return loc
}

func (m *ExportServiceImpl) getDates() []string {
	var dates []string
	m.DBMgr.GetDB().Model(&models.MeasurementResult{}).Distinct().Pluck("measurement_date", &dates)
	return dates
}

func addNewString(s []string, str string) []string {
	for _, v := range s {
		if v == str {
			return s
		}
	}

	return append(s, str)
}

func appendNewProbe(probes, probes2Append []*models.Probe) []*models.Probe {
	for _, m := range probes2Append {
		found := false
		for _, probe := range probes {
			if found = probe.ProbeID == m.ProbeID; found {
				break
			}
		}
		if !found {
			probes = append(probes, m)
		}
	}
	return probes
}
