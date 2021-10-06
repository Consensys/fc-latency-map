package locations

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ConsenSys/fc-latency-map/manager/constants"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type Airport struct {
	Name        string `json:"name"`
	Continent   string `json:"continent"`
	Coordinates string `json:"coordinates"`
	IataCode    string `json:"iata_code"`
	IsoCountry  string `json:"iso_country"`
	Type        string `json:"type"`
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
		log.Printf("ID:%d - Iata code: %s - Country code: %s - Name: %s - Type: %s\n",
			location.ID,
			location.IataCode,
			location.Country,
			location.Name,
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

func (srv *LocationServiceImpl) UpdateLocations(airportType, filename string) error {
	var airportTypeFormated string
	switch airportType {
	case constants.AirportTypeLarge:
		airportTypeFormated = "large_airport"
	case constants.AirportTypeMedium:
		airportTypeFormated = "medium_airport"
	case constants.AirportTypeSmall:
		airportTypeFormated = "small_airport"
	default:
		return errors.New("airport type not found")
	}

	airports, err := srv.ExtractAirports(filename)
	if err != nil {
		return err
	}

	cpt := 0
	for _, airport := range airports {
		if airport.Type == airportTypeFormated {
			existsLocation := models.Location{}
			srv.DBMgr.GetDB().Where("iata_code = ?", airport.IataCode).First(&existsLocation)
			if existsLocation == (models.Location{}) {
				coords := strings.Split(airport.Coordinates, ", ")
				lat, _ := strconv.ParseFloat(coords[1], 64)
				long, _ := strconv.ParseFloat(coords[0], 64)


				

				srv.DBMgr.GetDB().Create(&models.Location{
					Name:      airport.Name,
					Country:   airport.IsoCountry,
					IataCode:  getIataCode(airport.IataCode, airport.Name),
					Latitude:  lat,
					Longitude: long,
					Type:      airport.Type,
				})
				cpt++
			}
		}
	}
	log.Printf("%d airport imported, type: %s\n", cpt, airportTypeFormated)

	return nil
}


func formatName(name string) string {
	name = strings.Replace(name, "\u00e2\u0080\u0093", "-", -1)

	name = strings.Replace(name, "\u00c3\u0087", "Ç", -1)	
	name = strings.Replace(name, "\u00c3\u00a7", "ç", -1)

	name = strings.Replace(name, "\u00c3\u00ad", "í", -1)

	name = strings.Replace(name, "\u00c3\u00ba", "ú", -1)

	name = strings.Replace(name, "\u00c3\u00b3", "ó", -1)
	name = strings.Replace(name, "\u00c3\u0093", "Ó", -1)
	name = strings.Replace(name, "\u00c3\u00b4", "ô", -1)
	name = strings.Replace(name, "\u00c3\u00b6", "ö", -1)
	name = strings.Replace(name, "\u00c3\u00b8", "ø", -1)

	name = strings.Replace(name, "\u00c5\u0084", "ń", -1)
	name = strings.Replace(name, "\u00c3\u00b1", "ñ", -1)

	name = strings.Replace(name, "\u00c5\u0082", "ł", -1)
	
	name = strings.Replace(name, "\u00c5\u0084", "å", -1)
	name = strings.Replace(name, "\u00c5\u0081", "Å", -1)
	name = strings.Replace(name, "\u00c3\u00a1", "á", -1)
	name = strings.Replace(name, "\u00c3\u00a3", "ã", -1)
	name = strings.Replace(name, "\u00c3\u00a0", "à", -1)
	name = strings.Replace(name, "\u00c4\u0083", "ă", -1)

	name = strings.Replace(name, "\u00c4\u0099", "ę", -1)
	name = strings.Replace(name, "\u00c3\u00a9", "é", -1)
	name = strings.Replace(name, "\u00c3\u00a8", "è", -1)

	name = strings.Replace(name, "\u00c3\u009c", "Ü", -1)
	name = strings.Replace(name, "\u00c3\u00bc", "ü", -1)

	name = strings.Replace(name, "\u00c5\u00a1", "š", -1)
	name = strings.Replace(name, "\u00c5\u00a0", "Š", -1)

	name = strings.Replace(name, "\u00c5\u00be", "ž", -1)

	name = strings.Replace(name, "\u00c4\u008d", "č", -1)

	return name
}

func getIataCode(iata, name string) string {
	if iata == "" {
		return name
	}
	return iata
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

func (srv *LocationServiceImpl) FindAirport(iataCode, filename string) (Airport, error) {
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
