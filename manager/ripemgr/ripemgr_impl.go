package ripemgr

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/addresses"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	atlas "github.com/keltia/ripe-atlas"
)

type ServiceImpl struct {
	conf *viper.Viper
	ripe *atlas.Client
}

func NewServiceImpl(conf *viper.Viper, ripe *atlas.Client) RipeService {
	return &ServiceImpl{
		conf: conf,
		ripe: ripe,
	}
}

func (m *ServiceImpl) getMeasurementResults(ms map[int]int) ([]atlas.MeasurementResult, error) {
	var results []atlas.MeasurementResult
	for k, v := range ms {
		m.ripe.SetOption("start", strconv.Itoa(v))
		measurementResult, err := m.ripe.GetResults(k)
		if err != nil {
			return nil, err
		}

		results = append(results, measurementResult.Results...)
	}

	return results, nil
}

func (m *ServiceImpl) createMeasurements(miners []*models.Miner, probeIDs string) ([]*atlas.Measurement, error) {
	probes := []atlas.ProbeSet{
		{
			Type:      "probes",
			Value:     probeIDs,
			Requested: viper.GetInt("RIPE_REQUESTED_PROBES"),
		},
	}

	return m.createPing(miners, probes)
}

func (m *ServiceImpl) createPing(miners []*models.Miner, probes []atlas.ProbeSet) ([]*atlas.Measurement, error) {
	var d []atlas.Definition

	pingInterval := m.conf.GetInt("RIPE_PING_INTERVAL")

	for _, miner := range miners {
		if miner.IP == "" {
			continue
		}
		for _, ip := range strings.Split(miner.IP, ",") {
			ipAdd := net.ParseIP(ip)
			if ipAdd.IsPrivate() {
				continue
			}

			d = append(d, atlas.Definition{
				Description: fmt.Sprintf("%s ping to %s", miner.Address, ip),
				AF:          addresses.GetIPVersion(ipAdd),
				Target:      ip,
				Tags: []string{
					miner.Address,
				},
				Type:     "ping",
				Interval: pingInterval,
			})
		}
	}

	isOneOff := m.conf.GetBool("RIPE_ONE_OFF")
	runningTime := m.conf.GetInt("RIPE_PING_RUNNING_TIME")

	mr := &atlas.MeasurementRequest{
		Definitions: d,
		StartTime:   int(time.Now().Unix()),
		StopTime:    int(time.Now().Unix()) + runningTime,
		IsOneoff:    isOneOff,
		Probes:      probes,
	}

	p, err := m.ripe.Ping(mr)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
			"msg": mr,
		}).Error("Create ping")

		return nil, err
	}

	log.WithFields(log.Fields{
		"id":           p,
		"isOneOff":     isOneOff,
		"pingInterval": pingInterval,
		"measurement":  fmt.Sprintf("%#v\n", d),
	}).Info("creat newMeasurement")

	var measurement []*atlas.Measurement
	for _, v := range p.Measurements {
		measurement = append(measurement, &atlas.Measurement{
			ID:        v,
			StopTime:  mr.StopTime,
			StartTime: mr.StartTime,
			IsOneoff:  mr.IsOneoff,
		})
	}

	return measurement, err
}
