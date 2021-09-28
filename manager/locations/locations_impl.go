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
	Continent   string `json:"continent"`
	Coordinates string `json:"coordinates"`
	IataCode    string `json:"iata_code"`
	IsoCountry  string `json:"iso_country"`
}

type LocationServiceImpl struct {
	Conf  *viper.Viper
	DBMgr db.DatabaseMgr
}

func NewLocationServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr) LocationService {
	return &LocationServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
	}
}

func (srv *LocationServiceImpl) DisplayLocations() []*models.Location {
	locsList := []*models.Location{}
	srv.DBMgr.GetDB().Find(&locsList)
	for _, location := range locsList {
		log.Printf("ID:%d - Iata code: %s - Country code: %s\n", location.ID, location.IataCode, location.Country)
	}
	return locsList
}

func (srv *LocationServiceImpl) GetLocation(location *models.Location) *models.Location {
	if err := srv.DBMgr.GetDB().Where(location).First(&location).Error; err != nil {
		return nil
	}
	return location
}

func (srv *LocationServiceImpl) AddLocation(newLocation *models.Location) *models.Location {
	location := models.Location{}
	srv.DBMgr.GetDB().Where("iata_code = ?", newLocation.IataCode).First(&location)
	if location == (models.Location{}) {
		srv.DBMgr.GetDB().Create(&newLocation)
		log.Printf("New location, ID:%d - Iata code: %s - Country code: %s\n", newLocation.ID, newLocation.IataCode, newLocation.Country)
	} else {
		log.Printf("Location already exists, ID:%d\n", location.ID)
	}
	return &location
}

func (srv *LocationServiceImpl) DeleteLocation(location *models.Location) bool {
	l := srv.GetLocation(location)
	if l == nil {
		log.Printf("Unable to find location %s\n", location.Country)
	} else {
		srv.DBMgr.GetDB().Delete(&location)
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

	airports := make([]Airport, 0)
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
