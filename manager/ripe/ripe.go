package ripe

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	atlas "github.com/keltia/ripe-atlas"
)

var apiconfig = config.Config()

func GetProbes(countryCode string) ([]atlas.Probe, error) {	

	apiKey := apiconfig.GetString("RIPE_API_KEY")
	config :=atlas.Config{
		APIKey: apiKey,
	}

	client, err := atlas.NewClient(config)
	if err != nil {
		return nil, err
	}

	opts := make(map[string]string)
	opts["country_code"] = countryCode

	probes, err := client.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil
}