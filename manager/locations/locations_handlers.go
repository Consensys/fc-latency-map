package locations

import (
	"strconv"
	"strings"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationHandler struct {
	LSer LocationService
}

func NewLocationHandler() *LocationHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	lSer := NewLocationServiceImpl(conf, dbMgr)
	return &LocationHandler{
		LSer: lSer,
	}
}

// GetLocations handle locations get cli command
func (mHdl *LocationHandler) GetLocations() { //nolint:revive
	mHdl.LSer.DisplayLocations()
}

// AddLocation handle location add cli command
func (mHdl *LocationHandler) AddLocation(airportCode string) (*models.Location, error) {
	airport, err := mHdl.LSer.FindAirport(airportCode)
	if err != nil {
		return nil, err
	}

	coords := strings.Split(airport.Coordinates, ", ")
	lat, _ := strconv.ParseFloat(coords[1], 32)
	long, _ := strconv.ParseFloat(coords[0], 32)
	location := &models.Location{
		Country:   airport.IsoCountry,
		IataCode:  airport.IataCode,
		Latitude:  lat,
		Longitude: long,
	}
	location = mHdl.LSer.AddLocation(location)

	return location, nil
}

// DeleteLocation handle location delete cli command
func (mHdl *LocationHandler) DeleteLocation(countryCode string) {
	location := &models.Location{
		Country: countryCode,
	}
	mHdl.LSer.DeleteLocation(location)
}
