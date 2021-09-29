package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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
	Type  			string `json:"type"`
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

func (srv *LocationServiceImpl) GetAllLocations() []*models.Location {
	locsList := []*models.Location{}
	srv.DBMgr.GetDB().Find(&locsList)
	for _, location := range locsList {
		log.Printf("ID:%d - Iata code: %s - Country code: %s - Type: %s\n", 
		location.ID,
		location.IataCode,
		location.Country,
		location.Type)
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
		log.Printf("new location, ID:%d - Country code: %s\n", newLocation.ID, newLocation.Country)
	} else {
		log.Printf("location already exists, ID:%d\n", location.ID)
	}
	return &location
}

func (srv *LocationServiceImpl) DeleteLocation(location *models.Location) bool {
	if l := srv.GetLocation(location); l == nil {
		log.Printf("unable to find location %s\n", location.Country)
	} else {
		srv.DBMgr.GetDB().Delete(&location)
		log.Printf("location %d deleted\n", location.ID)
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

func (srv *LocationServiceImpl) UpdateLocations(airportType string, filename string) error {
	var airportTypeFormated string
	switch airportType {
		case "large":
				airportTypeFormated = "large_airport"
		case "medium":
				fmt.Println("It's the medium weekend")
				airportTypeFormated = "medium_airport"
		case "small":
				fmt.Println("It's the small weekend")
				airportTypeFormated = "small_airport"
		default:
			return errors.New("airport type not found")
	}
	log.Printf("import airport type: %s\n", airportTypeFormated)

	airports, err := srv.ExtractAirports(filename)
	if err != nil {
		return nil
	}

	cpt := 0
	for _, airport := range airports {
		if airport.Type == airportTypeFormated {
			coords := strings.Split(airport.Coordinates, ", ")
			lat, _ := strconv.ParseFloat(coords[0], 64)
			long, _ := strconv.ParseFloat(coords[1], 64)
			location := &models.Location{
				Country:   airport.IsoCountry,
				IataCode:  airport.IataCode,
				Latitude:  lat,
				Longitude: long,
				Type: airport.Type,
			}

			existsLocation := models.Location{}
			srv.DBMgr.GetDB().Where("iata_code = ?", location.IataCode).First(&existsLocation)
			if existsLocation == (models.Location{}) {
				srv.DBMgr.GetDB().Create(&location)
				cpt++
			}
		}		
	}
	log.Printf("%d airport imported\n", cpt)

	return nil
}

func (srv *LocationServiceImpl) ExtractAirports(filename string) ([]Airport, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return []Airport{}, err
	}
	log.Printf("successfully Opened: %s\n", filename)
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []Airport{}, err
	}

	airports := make([]Airport, 0)
	if err := json.Unmarshal(jsonData, &airports); err != nil {
		return []Airport{}, err
	}

	return airports, nil
}

func (srv *LocationServiceImpl) FindAirport(iataCode string, filename string) (Airport, error) {
	airports, err := srv.ExtractAirports(filename)
	if err != nil {
		return Airport{}, err
	}

	for _, airport := range airports {
		if airport.IataCode == iataCode {
			return airport, nil
		}
	}

	return Airport{}, errors.New("airport IataCode not found")
}
