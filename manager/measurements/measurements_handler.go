package measurements

import (
	"encoding/json"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

func Export(fn string) {
	measurements := GetLatencyMeasurements()

	fullJson, err := json.MarshalIndent(measurements, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create json data")
		return
	}

	file.Create(fn, fullJson)
	log.WithFields(log.Fields{
		"file": fn,
	}).Info("Export successful")
}

func GetLatencyMeasurements() []*models.Latency {
	var latencies = []*models.Latency{}
	for m := 0; m < 10; m++ {

		latency := &models.Latency{
			Address: randomString(10),
		}
		latencies = append(latencies, latency)

		for l := 0; l < 3; l++ {

			location := &models.LocationData{
				Country:   randomString(10),
				Latitude:  randomString(10),
				Longitude: randomString(10),
			}
			latency.Locations = append(latency.Locations, location)
			for meas := 0; meas < 5; meas++ {
				measureData := &models.MeasureData{
					Avg: rand.Float64(),
					Lts: rand.Int(),
					Max: rand.Float64(),
					Min: rand.Float64(),
				}
				location.Measures = append(location.Measures, measureData)
			}
		}

	}
	return latencies
}

// ///////////////////////////////////////////////////
func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
