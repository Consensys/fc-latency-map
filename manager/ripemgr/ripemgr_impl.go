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

type RipeMgrImpl struct {
	conf *viper.Viper
	c    *atlas.Client
}

func NewRipeImpl(conf *viper.Viper) (RipeMgr, error) {
	cfgs := []atlas.Config{}
	cfgs = append(cfgs, atlas.Config{
		APIKey: conf.GetString("RIPE_API_KEY"),
	})
	c, err := atlas.NewClient(cfgs...)
	if err != nil {
		log.Println("Connecting to Ripe Atlas API", err)
		return nil, err
	}
	ver := atlas.GetVersion()
	log.Println("api version ", ver)

	return &RipeMgrImpl{
		c: c,
	}, nil
}

func (fMgr *RipeMgrImpl) GetProbe(id int) (*atlas.Probe, error) {
	return fMgr.c.GetProbe(id)
}

func (fMgr *RipeMgrImpl) GetProbes(opts map[string]string) ([]atlas.Probe, error) {
	probes, err := fMgr.c.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil

}

func (m *RipeMgrImpl) GetMeasurementResults(ms map[int]int) ([]atlas.MeasurementResult, error) {
	var results []atlas.MeasurementResult
	for k, v := range ms {
		m.c.SetOption("start", strconv.Itoa(v))
		measurementResult, err := m.c.GetResults(k)
		if err != nil {
			return nil, err
		}

		results = append(results, measurementResult.Results...)
	}

	return results, nil
}

func (m *RipeMgrImpl) CreateMeasurements(miners []*models.Miner, probeIDs string) ([]*atlas.Measurement, error) {
	probes := []atlas.ProbeSet{
		{
			Type:      "probes",
			Value:     probeIDs,
			Requested: viper.GetInt("RIPE_REQUESTED_PROBES"),
		},
	}

	return m.createPing(miners, probes)
}

func (m *RipeMgrImpl) createPing(miners []*models.Miner, probes []atlas.ProbeSet) ([]*atlas.Measurement, error) {
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

	p, err := m.c.Ping(mr)

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
