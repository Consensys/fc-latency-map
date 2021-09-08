package probes

import (
	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
)

type Ripe struct {
	c         *atlas.Client
	StartTime int
	StopTime  int
	IsOneOff  bool
}

func NewClient(t string, cfgs ...atlas.Config) (*Ripe, error) {
	r := &Ripe{}
	if cfgs == nil {
		cfgs = append(cfgs, atlas.Config{
			APIKey: t,
		})
	}

	c, err := atlas.NewClient(cfgs...)
	if err != nil {
		log.Println("Connecting to Ripe Atlas API", err)
		return nil, err
	}
	r.c = c
	ver := atlas.GetVersion()
	log.Println("api version ", ver)

	return r, nil
}

func (r *Ripe) GetProbe(id int) (m *atlas.Probe, err error) {
	return r.c.GetProbe(id)
}

func (r *Ripe) GetProbes(countryCode string) ([]atlas.Probe, error) {
	opts := make(map[string]string)
	opts["country_code"] = countryCode

	probes, err := r.c.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil
}

func (r *Ripe) GetBestProbes(countryProbes []atlas.Probe) (atlas.Probe, error) {
	for _, probe := range countryProbes {
		if probe.Status.Name == "Connected" {
			return probe, nil
		}
	}
	return atlas.Probe{}, nil
}

func (r *Ripe) GetAllProbes() ([]atlas.Probe, error) {

	var countries [2]string
	countries[0] = "FR"
	countries[1] = "PT"

	var bestProbes []atlas.Probe

	for _, country := range countries {
		log.WithFields(log.Fields{
			"country": country,
		}).Info("Get probes for country")
		countryProbes, err := r.GetProbes(country)
		if err != nil {
			return nil, err
		}
		bestProbe, err := r.GetBestProbes(countryProbes)
		if err != nil {
			return nil, err
		}
		bestProbes = append(bestProbes, bestProbe)
	}
	return bestProbes, nil
}
<<<<<<< HEAD

func (r *Ripe) Update() {
	// get countries from db
	// countries := []string{"FR", "PT"}

	probes, err := r.GetAllProbes() // by countries
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetAllProbes")
		return
	}

	// update db probes

	for i, probe := range probes {
		log.Debug(i, probe)
	}
}
=======
>>>>>>> main
