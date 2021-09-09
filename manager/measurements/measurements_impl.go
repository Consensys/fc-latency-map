package measurements

import (
	"fmt"
	"time"

	"github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Ripe struct {
	c *atlas.Client
}

func NewClient(apiKey string, cfgs ...atlas.Config) (*Ripe, error) {
	r := &Ripe{}
	if cfgs == nil {
		cfgs = append(cfgs, atlas.Config{
			APIKey: apiKey,
		})
	}

	c, err := atlas.NewClient(cfgs...)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("connection to api")
		return nil, err
	}
	r.c = c

	log.WithFields(log.Fields{
		"version": atlas.GetVersion(),
	}).Info("get api version")

	return r, nil
}

func (r *Ripe) GetMeasurement(id int) (m *atlas.Measurement, err error) {
	return r.c.GetMeasurement(id)
}

func (r *Ripe) CreatePing(miners []Miner, probes []atlas.ProbeSet) (*atlas.MeasurementResp, error) {
	var d []atlas.Definition

	pingInterval := viper.GetInt("RIPE_PING_INTERVAL")

	for _, miner := range miners {
		for _, ip := range miner.Ip {
			d = append(d, atlas.Definition{
				Description: fmt.Sprintf("%s ping to %s", miner.Address, ip),
				AF:          4,
				Target:      ip,
				Tags: []string{
					miner.Address,
				},
				Type:     "ping", // 10 minutes
				Interval: pingInterval,
			})
		}
	}

	isOneOff := viper.GetBool("RIPE_ONE_OFF")

	mr := &atlas.MeasurementRequest{
		Definitions: d,
		StartTime:   int(time.Now().Unix()),
		StopTime:    int(time.Now().Unix() + 3600), // 1 hour
		IsOneoff:    isOneOff,
		Probes:      probes,
	}

	p, err := r.c.Ping(mr)
	log.WithFields(log.Fields{
		"id":           p,
		"isOneOff":     isOneOff,
		"pingInterval": pingInterval,
		"measurement":  fmt.Sprintf("%#v\n", d),
	}).Info("creat newMeasurement")

	return p, err
}

func (r *Ripe) GetMeasurementResult(id int) (probeResults []MeasurementResults, err error) {

	var probes map[int]*atlas.Probe
	results, err := r.c.GetResults(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Info("get results")
		return nil, err
	}

	// create a str
	for _, ripeResult := range results.Results {
		p, found := probes[ripeResult.PrbID]
		if !found {
			p, err = r.c.GetProbe(ripeResult.PrbID)
			if err != nil {
				log.WithFields(log.Fields{
					"id":  ripeResult.PrbID,
					"err": err,
				}).Info("get probe")
				continue
			}
		}

		probeResults = append(probeResults, MeasurementResults{
			Measurement: ripeResult,
			Probe:       *p,
		})
	}
	return probeResults, err
}
