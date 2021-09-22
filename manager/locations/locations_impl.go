package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/constants"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type Airport struct {
	Continent 	string `json:"continent"`
	Coordinates string `json:"coordinates"`
	IataCode 		string `json:"iata_code"`
	IsoCountry 	string `json:"iso_country"`
}

type LocationServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
}

func NewLocationServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr) LocationService {
	return &LocationServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
	}
}

func (srv *LocationServiceImpl) GetLocations() []*models.Location {
	var locsList = []*models.Location{}
	(*srv.DbMgr).GetDb().Find(&locsList)
	for _, location := range locsList {
		log.Printf("ID:%d - Country code: %s\n", location.ID, location.Country)
	}
	return locsList
}

func (srv *LocationServiceImpl) GetLocation(location models.Location) models.Location {
	if err := (*srv.DbMgr).GetDb().Where(location).First(&location).Error; err != nil {
		return models.Location{}
	}
	return location
}

func (srv *LocationServiceImpl) AddLocation(newLocation models.Location) models.Location {
	var location = models.Location{}
	(*srv.DbMgr).GetDb().Where("iata_code = ?", newLocation.IataCode).First(&location)
	if location == (models.Location{}) {
		(*srv.DbMgr).GetDb().Create(&newLocation) 
		log.Printf("New location, ID:%d - Country code: %s\n", newLocation.ID, newLocation.Country)
	} else {
		log.Printf("Location already exists, ID:%d\n", location.ID)
	}
	return location
}

func (srv *LocationServiceImpl) DeleteLocation(location models.Location) bool {
	location = srv.GetLocation(location)
	if (location == models.Location{}) {
		log.Printf("Unable to find location %s\n", location.Country)
	} else {
		(*srv.DbMgr).GetDb().Delete(&location)
		log.Printf("Location %d deleted\n", location.ID)
	}
	
	return true
}

func (srv *LocationServiceImpl) CheckCountry(countryCode string) bool {	
	for _, country := range constants.Countries {
		if countryCode == country.Code {
			return true
		}
	}
	return false
}

func (srv *LocationServiceImpl) FindAirport(iataCode string) (Airport, error) {
	filename := "constants/airport-codes.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		return Airport{}, err
	}
	fmt.Printf("Successfully Opened: %s\n", filename)
	defer jsonFile.Close()
	
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Airport{}, err
	}

	airports := make([]Airport,0)
	if err := json.Unmarshal(jsonData, &airports); err != nil {
		return Airport{}, err
	}

	for _, airport := range airports {
		if airport.IataCode == iataCode {
			return airport, nil
		}
	}
	
	return Airport{}, errors.New("airport IataCode not found")
}