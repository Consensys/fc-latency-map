package measurements

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

func (m *MeasurementServiceImpl) GetRipeMeasurement(id int) (*atlas.Measurement, error) {
	return m.Ripe.GetMeasurement(id)
}

func (m *MeasurementServiceImpl) CreatePingProbes(miners []*models.Miner, t, value string) (*atlas.MeasurementResp, error) {
	probes := []atlas.ProbeSet{
		{
			Type:      t,
			Value:     value,
			Requested: viper.GetInt("RIPE_REQUESTED_PROBES"),
		},
	}
	return m.CreatePing(miners, probes)
}

func (m *MeasurementServiceImpl) CreatePing(miners []*models.Miner, probes []atlas.ProbeSet) (*atlas.MeasurementResp, error) {
	var d []atlas.Definition

	pingInterval := m.Conf.GetInt("RIPE_PING_INTERVAL")

	for _, miner := range miners {
		for _, ip := range strings.Split(miner.Ip, ",") {
			d = append(d, atlas.Definition{
				Description: fmt.Sprintf("%s ping to %s", miner.Address, ip),
				AF:          4,
				Target:      ip,
				Tags: []string{
					miner.Address,
				},
				Type:     "ping",
				Interval: pingInterval,
			})
		}
	}

	isOneOff := m.Conf.GetBool("RIPE_ONE_OFF")

	mr := &atlas.MeasurementRequest{
		Definitions: d,
		StartTime:   int(time.Now().Unix()),
		StopTime:    int(time.Now().Unix() + 3600), // 1 hour
		IsOneoff:    isOneOff,
		Probes:      probes,
	}

	p, err := m.Ripe.Ping(mr)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
			"msg": mr,
		}).Info("Create ping")
		return nil, err
	}
	log.WithFields(log.Fields{
		"id":           p,
		"isOneOff":     isOneOff,
		"pingInterval": pingInterval,
		"measurement":  fmt.Sprintf("%#v\n", d),
	}).Info("creat newMeasurement")

	return p, err
}

func (m *MeasurementServiceImpl) GetRipeMeasurementResult(id int) ([]atlas.MeasurementResult, error) {

	results, err := m.Ripe.GetResults(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Info("get results")
		return nil, err
	}
	return results.Results, err
}

func (m *MeasurementServiceImpl) GetRipeMeasures(tag string) {

	measurementResults, err := m.getRipeMeasurementResults(tag)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Load measurement Results from Ripe")
	}

	m.importMeasurement(measurementResults)

	log.Info("measurements successfully get")
}

func (m *MeasurementServiceImpl) getRipeMeasurementResults(tag string) ([]MeasurementResult, error) {
	ops := make(map[string]string)
	ops["tags"] = tag

	measurements, err := m.Ripe.GetMeasurements(ops)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("GetRipeMeasurementResult")
		return nil, err
	}
	var measurementResults []MeasurementResult
	for _, measurement := range measurements {
		results, err := m.GetRipeMeasurementResult(measurement.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"id":  measurement.ID,
				"err": err,
			}).Info("GetRipeMeasurementResult")
			continue
		}
		measurementResults = append(measurementResults, MeasurementResult{
			Measurement: measurement,
			Results:     results,
		})
	}
	return measurementResults, err
}
