package locations

import (
	"fmt"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

// GetLocationsHandler handle locations get cli command
func GetLocationsHandler() {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	locs := NewLocationServiceImpl(dbMgr)
		locsList := locs.GetLocations()
		for _, location := range locsList {
			fmt.Printf("ID:%d - Country code: %s\n", location.ID, location.Country)
		}
}

// AddLocationHandler handle location add cli command
func AddLocationHandler(countryCode string) {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	newLocation := models.Location{
		Country: countryCode,
		Latitude:    "1.2",
		Longitude: "2.1",
	}
	var location = models.Location{}
	// dbMgr.GetDb().Where(&newLocation).First(&location)
	dbMgr.GetDb().Where(&newLocation).First(&location)
	fmt.Printf("xxx> %s\n", location)
	if location == (models.Location{}) {
		locs := NewLocationServiceImpl(dbMgr)
		newLocation = locs.AddLocation(newLocation)
		fmt.Printf("New location, ID:%d - Country code: %s\n", newLocation.ID, newLocation.Country)
	} else {
		fmt.Printf("Location already exists, ID:%d\n", location.ID)
	}
	
}


// DeleteLocationHandler handle location delete cli command
func DeleteLocationHandler(countryCode string) {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	location := models.Location{
		Country: countryCode,
	}
	locs := NewLocationServiceImpl(dbMgr)
	location = locs.GetLocation(location)
	if (location == models.Location{}) {
		fmt.Printf("Unable to find location %s\n", countryCode)
	} else {
		locs.DeleteLocation(location)
		fmt.Printf("Location %d deleted\n", location.ID)
	}
}