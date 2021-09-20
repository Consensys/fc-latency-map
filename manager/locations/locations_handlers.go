package locations

import (
	"errors"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/constants"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationHandler struct {
	LSer *LocationService
}

func NewLocationHandler() *LocationHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	lSer := NewLocationServiceImpl(conf, &dbMgr)
	return &LocationHandler{
		LSer: &lSer,
	}
}

// GetLocationsHandler handle locations get cli command
func (mHdl *LocationHandler) GetLocations() {
	(*mHdl.LSer).GetLocations()
}

// AddLocationHandler handle location add cli command
func (mHdl *LocationHandler) AddLocation(countryCode string)  (models.Location, error) {
	if !checkCountry(countryCode) {
		err := errors.New("country code not found")
		return models.Location{}, err
	}

	location := models.Location{
		Country: countryCode,
		Latitude:    "0",
		Longitude: "0",
	}
	location = (*mHdl.LSer).AddLocation(location)
	return location, nil
}

// DeleteLocation handle location delete cli command
func (mHdl *LocationHandler) DeleteLocation(countryCode string) {
	location := models.Location{
		Country: countryCode,
	}
	(*mHdl.LSer).DeleteLocation(location)
}

// checkCountry checks the country code exists
func checkCountry(countryCode string) bool {
	for _, country := range constants.Countries {
		if countryCode == country.Code {
			return true
		}
	}
	return false
}