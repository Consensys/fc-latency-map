package locations

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationHandler struct {
	MSer *LocationService
}

func NewLocationHandler() *LocationHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	mSer := NewLocationServiceImpl(conf, &dbMgr)
	return &LocationHandler{
		MSer: &mSer,
	}
}

// GetLocationsHandler handle locations get cli command
func (mHdl *LocationHandler) GetLocations() {
	(*mHdl.MSer).GetLocations()
}

// AddLocationHandler handle location add cli command
func (mHdl *LocationHandler) AddLocation(countryCode string) {
	location := models.Location{
		Country: countryCode,
		Latitude:    "0",
		Longitude: "0",
	}
	(*mHdl.MSer).AddLocation(location)
}

// DeleteLocation handle location delete cli command
func (mHdl *LocationHandler) DeleteLocation(countryCode string) {
	location := models.Location{
		Country: countryCode,
	}
	(*mHdl.MSer).DeleteLocation(location)
}