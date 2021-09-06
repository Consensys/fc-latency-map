package probes

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ConsenSys/fc-latency-map/manager/ripe"
	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
)

// GetProbes returns probes list for measurement
func GetProbes() ([]atlas.Probe, error) {
	jsonFile, err := os.Open("constants/countries.json")
	if err != nil {
    return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var countries []string
	var bestProbes []atlas.Probe
	json.Unmarshal(byteValue, &countries)
	for _, country := range countries {
		log.WithFields(log.Fields{
			"country": country,
		}).Info("Get probes for country")
		countryProbes, err := ripe.GetProbes(country)
		if err != nil {
			return nil, err
		}
		bestProbe, err := GetBestProbes(countryProbes)
		if err != nil {
			return nil, err
		}
		bestProbes = append(bestProbes, bestProbe)
	}
	return bestProbes, nil
}

// GetBestProbes returns selected online probes
func GetBestProbes(countryProbes []atlas.Probe) (atlas.Probe, error) {
	for _, probe := range countryProbes {
		if probe.Status.Name == "Connected" {
			return probe, nil
		}
	}
	return atlas.Probe{}, nil
}