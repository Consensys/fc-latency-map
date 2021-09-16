package measurements

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

func (m *MeasurementServiceImpl) RipeGetMeasurement(id int) (*atlas.Measurement, error) {
	return m.Ripe.GetMeasurement(id)
}

func (m *MeasurementServiceImpl) RipeCreatePingWithProbes(miners []*models.Miner, probeIDs string) (*atlas.MeasurementRequest, *atlas.MeasurementResp, error) {
	probes := []atlas.ProbeSet{
		{
			Type:      "probes",
			Value:     probeIDs,
			Requested: viper.GetInt("RIPE_REQUESTED_PROBES"),
		},
	}
	return m.RipeCreatePing(miners, probes)

}

func (m *MeasurementServiceImpl) RipeCreatePing(miners []*models.Miner, probes []atlas.ProbeSet) (*atlas.MeasurementRequest, *atlas.MeasurementResp, error) {
	var d []atlas.Definition

	pingInterval := m.Conf.GetInt("RIPE_PING_INTERVAL")

	for _, miner := range miners {
		for _, ip := range strings.Split(miner.Ip, ",") {
			if net.ParseIP(ip).IsPrivate() {
				continue
			}
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
	runningTime := m.Conf.GetInt("RIPE_PING_RUNNING_TIME")

	mr := &atlas.MeasurementRequest{
		Definitions: d,
		StartTime:   int(time.Now().Unix()),
		StopTime:    int(time.Now().Unix()) + runningTime,
		IsOneoff:    isOneOff,
		Probes:      probes,
	}

	p, err := m.Ripe.Ping(mr)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
			"msg": mr,
		}).Info("Create ping")

		return nil, nil, err
	}

	log.WithFields(log.Fields{
		"id":           p,
		"isOneOff":     isOneOff,
		"pingInterval": pingInterval,
		"measurement":  fmt.Sprintf("%#v\n", d),
	}).Info("creat newMeasurement")

	return mr, p, err
}

func (m *MeasurementServiceImpl) RipeGetMeasurementResult(id int, start int) ([]atlas.MeasurementResult, error) {
	m.Ripe.SetOption("start", strconv.Itoa(start))
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

func (m *MeasurementServiceImpl) getRipeMeasurementResultsById(id int, start int) ([]MeasurementResult, error) {
	var measurementResults []MeasurementResult
	measurement, _ := m.Ripe.GetMeasurement(id)

	results, err := m.RipeGetMeasurementResult(id, start)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Info("RipeGetMeasurementResult")

	}
	measurementResults = append(measurementResults, MeasurementResult{
		Measurement: *measurement,
		Results:     results,
	})

	return measurementResults, err
}

func (m *MeasurementServiceImpl) getRipeMeasurementResultsByTag(tag string) ([]MeasurementResult, error) {
	ops := make(map[string]string)
	ops["tags"] = tag

	measurements, err := m.Ripe.GetMeasurements(ops)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("RipeGetMeasurementResult")
		return nil, err
	}
	var measurementResults []MeasurementResult
	for _, measurement := range measurements {
		results, err := m.RipeGetMeasurementResult(measurement.ID, 0)
		if err != nil {
			log.WithFields(log.Fields{
				"id":  measurement.ID,
				"err": err,
			}).Info("RipeGetMeasurementResult")
			continue
		}
		measurementResults = append(measurementResults, MeasurementResult{
			Measurement: measurement,
			Results:     results,
		})
	}
	return measurementResults, err
}
