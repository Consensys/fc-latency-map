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
		c:    c,
		conf: conf,
	}, nil
}

func (rMge *RipeMgrImpl) GetProbe(id int) (*atlas.Probe, error) {
	return rMge.c.GetProbe(id)
}

func (rMge *RipeMgrImpl) GetProbes(opts map[string]string) ([]atlas.Probe, error) {
	probes, err := rMge.c.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil

}

func (rMge *RipeMgrImpl) GetNearestProbe(latitude, longitude string) (*atlas.Probe, error) {
	var err error
	nearestProbes := []atlas.Probe{}

	maxLocRange := rMge.conf.GetFloat64("RIPE_LOCATION_RANGE_MAX")
	coordRange := rMge.conf.GetFloat64("RIPE_LOCATION_RANGE_INIT")

	lat, _ := strconv.ParseFloat(latitude, 32)
	long, _ := strconv.ParseFloat(longitude, 32)

	for len(nearestProbes) < 1 && coordRange < maxLocRange {
		latGte := fmt.Sprintf("%f", lat-coordRange)
		latLte := fmt.Sprintf("%f", lat+coordRange)
		longGte := fmt.Sprintf("%f", long-coordRange)
		longLte := fmt.Sprintf("%f", long+coordRange)

		log.WithFields(log.Fields{
			"latitude__gte":  latGte,
			"latitude__lte":  latLte,
			"longitude__gte": longGte,
			"longitude__lte": longLte,
			"range":          coordRange,
		}).Info("Get probes for geo location")

		opts := make(map[string]string)
		opts["latitude__gte"] = latGte
		opts["latitude__lte"] = latLte
		opts["longitude__gte"] = longGte
		opts["longitude__lte"] = longLte
		opts["status_name"] = "Connected"
		opts["sort"] = "id"

		nearestProbes, err = rMge.c.GetProbes(opts)

		if err != nil {
			if err.Error() == "empty probe list" {
				coordRange = coordRange * 2
				continue
			}
			return nil, err
		}
	}
	return &nearestProbes[0], nil
}

func (rMge *RipeMgrImpl) GetMeasurementResults(ms map[int]int) ([]atlas.MeasurementResult, error) {
	var results []atlas.MeasurementResult
	for k, v := range ms {
		rMge.c.SetOption("start", strconv.Itoa(v))
		measurementResult, err := rMge.c.GetResults(k)
		if err != nil {
			return nil, err
		}

		results = append(results, measurementResult.Results...)
	}

	return results, nil
}

func (rMge *RipeMgrImpl) CreateMeasurements(miners []*models.Miner, probeIDs string) ([]*atlas.Measurement, error) {
	probes := []atlas.ProbeSet{
		{
			Type:      "probes",
			Value:     probeIDs,
			Requested: viper.GetInt("RIPE_REQUESTED_PROBES"),
		},
	}

	return rMge.createPing(miners, probes)
}

func (rMge *RipeMgrImpl) createPing(miners []*models.Miner, probes []atlas.ProbeSet) ([]*atlas.Measurement, error) {
	var d []atlas.Definition

	pingInterval := rMge.conf.GetInt("RIPE_PING_INTERVAL")

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

	isOneOff := rMge.conf.GetBool("RIPE_ONE_OFF")
	runningTime := rMge.conf.GetInt("RIPE_PING_RUNNING_TIME")

	mr := &atlas.MeasurementRequest{
		Definitions: d,
		StartTime:   int(time.Now().Unix()),
		StopTime:    int(time.Now().Unix()) + runningTime,
		IsOneoff:    isOneOff,
		Probes:      probes,
	}

	p, err := rMge.c.Ping(mr)

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
