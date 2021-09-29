package locations

import (
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationHandler struct {
	Conf  *viper.Viper
	LSer *LocationService
}

func NewLocationHandler() *LocationHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	lSer := NewLocationServiceImpl(conf, dbMgr)
	return &LocationHandler{
		Conf:  conf,
		LSer: &lSer,
	}
}

// GetAllLocations handle locations get cli command
func (mHdl *LocationHandler) GetAllLocations() {
	(*mHdl.LSer).GetAllLocations()
}

// UpdateLocations handle adding all airport in database
func (mHdl *LocationHandler) UpdateLocations(airportType string) error {
	err := (*mHdl.LSer).UpdateLocations(airportType, mHdl.Conf.GetString("CONSTANT_AIRPORTS"))
	return err
}

// AddLocation handle location add cli command
func (mHdl *LocationHandler) AddLocation(airportCode string) (*models.Location, error) {
	airport, err := (*mHdl.LSer).FindAirport(airportCode, mHdl.Conf.GetString("CONSTANT_AIRPORTS"))
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
	location = (*mHdl.LSer).AddLocation(location)

	return location, nil
}

// DeleteLocation handle location delete cli command
func (mHdl *LocationHandler) DeleteLocation(iataCode string) {
	location := &models.Location{
		IataCode: iataCode,
	}
	(*mHdl.LSer).DeleteLocation(location)
}
