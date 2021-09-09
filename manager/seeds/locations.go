package seeds

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

var locations = []models.Location{
	models.Location{
		Country: "FR",
		Latitude:    "1.2",
		Longitude: "2.1",
	},
	models.Location{
		Country: "PT",
		Latitude:    "3.5",
		Longitude: "6.7",
	},
}