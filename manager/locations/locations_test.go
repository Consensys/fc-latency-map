package locations

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

var dummyCountryError = "FRX"

var dummyName = "Charles de Gaulle International Airport"
var dummyCountry = "FR"
var dummyIataCode = "CDG"
var dummyGeoLatitude = 49.012798
var dummyGeoLongitude = 2.55
var dummyType = "large_airport"

var dummyLocation = models.Location{
	Name:   		dummyName,
	Country:   dummyCountry,
	IataCode:  dummyIataCode,
	Latitude:  dummyGeoLatitude,
	Longitude: dummyGeoLongitude,
	Type:      dummyType,
}

var mockAirportJson = `[{"continent": "EU", "coordinates": "1.90889, 48.843601", "elevation_ft": "371", "gps_code": "LFPF", "iata_code": null, "ident": "LFPF", "iso_country": "FR", "iso_region": "FR-IDF", "local_code": null, "municipality": "Creil", "name": "A\u00c3\u00a9rodrome de Beynes - Thiverval", "type": "small_airport"},{"continent": "EU", "coordinates": "2.55, 49.012798", "elevation_ft": "392", "gps_code": "LFPG", "iata_code": "CDG", "ident": "LFPG", "iso_country": "FR", "iso_region": "FR-IDF", "local_code": null, "municipality": "Paris", "name": "Charles de Gaulle International Airport", "type": "large_airport"},{"continent": "EU", "coordinates": "2.6075, 48.8978", "elevation_ft": "207", "gps_code": "LFPH", "iata_code": null, "ident": "LFPH", "iso_country": "FR", "iso_region": "FR-IDF", "local_code": null, "municipality": "Paris", "name": "A\u00c3\u00a9rodrome de Chelles-le-Pin", "type": "small_airport"}]`

func Test_GetAllLocations_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	locations := srv.GetAllLocations()

	// Assert
	assert.NotNil(t, locations)
	assert.Empty(t, locations)
}

func Test_GetAllLocations_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()

	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	srv.AddLocation(&dummyLocation)
	locations := srv.GetAllLocations()

	// Assert
	assert.NotNil(t, locations)
	assert.NotEmpty(t, locations)

	actual := *(locations[0])
	assert.Equal(t, dummyLocation.Name, actual.Name)
	assert.Equal(t, dummyLocation.Country, actual.Country)
	assert.Equal(t, dummyLocation.IataCode, actual.IataCode)
	assert.Equal(t, dummyLocation.Latitude, actual.Latitude)
	assert.Equal(t, dummyLocation.Longitude, actual.Longitude)
	assert.Equal(t, dummyLocation.Type, actual.Type)
}

func Test_GetTotalLocations_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()

	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	srv.AddLocation(&dummyLocation)
	count := srv.GetTotalLocations()

	// Assert
	assert.Equal(t, int64(1), count)
}

func Test_GetLocation_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	location := srv.GetLocation(&dummyLocation)

	// Assert
	assert.Empty(t, location)
}

func Test_GetLocation_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	srv.AddLocation(&dummyLocation)
	location := srv.GetLocation(&dummyLocation)

	// Assert
	assert.NotNil(t, location)
	assert.NotEmpty(t, location)

	assert.Equal(t, dummyLocation.Name, location.Name)
	assert.Equal(t, dummyLocation.Country, location.Country)
	assert.Equal(t, dummyLocation.IataCode, location.IataCode)
	assert.Equal(t, dummyLocation.Latitude, location.Latitude)
	assert.Equal(t, dummyLocation.Longitude, location.Longitude)
	assert.Equal(t, dummyLocation.Type, location.Type)
}

func Test_AddLocation_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	location := srv.GetLocation(&dummyLocation)

	// Assert
	assert.Empty(t, location)

	// Act
	srv.AddLocation(&dummyLocation)
	newLocation := srv.GetLocation(&dummyLocation)

	// Assert
	assert.NotNil(t, newLocation)
	assert.NotEmpty(t, newLocation)

	assert.Equal(t, dummyLocation.Name, newLocation.Name)
	assert.Equal(t, dummyLocation.Country, newLocation.Country)
	assert.Equal(t, dummyLocation.IataCode, newLocation.IataCode)
	assert.Equal(t, dummyLocation.Latitude, newLocation.Latitude)
	assert.Equal(t, dummyLocation.Longitude, newLocation.Longitude)
	assert.Equal(t, dummyLocation.Type, newLocation.Type)
}

func Test_DeleteLocation_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	srv.AddLocation(&dummyLocation)
	location := srv.GetLocation(&dummyLocation)

	// Assert
	assert.NotNil(t, location)
	assert.NotEmpty(t, location)

	// Act
	srv.DeleteLocation(&dummyLocation)
	deletedLocation := srv.GetLocation(&dummyLocation)

	// Assert
	assert.Empty(t, deletedLocation)
}

func Test_UpdateLocations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock file
	file, err := ioutil.TempFile("", "airports")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(mockAirportJson)

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	err = srv.UpdateLocations("large", file.Name())
	location := &models.Location{
		IataCode: dummyIataCode,
	}
	newLocation := srv.GetLocation(location)

	assert.Nil(t, err)
	assert.Equal(t, dummyLocation.Name, newLocation.Name)
	assert.Equal(t, dummyLocation.Country, newLocation.Country)
	assert.Equal(t, dummyLocation.IataCode, newLocation.IataCode)
	assert.Equal(t, dummyLocation.Latitude, newLocation.Latitude)
	assert.Equal(t, dummyLocation.Longitude, newLocation.Longitude)
	assert.Equal(t, dummyLocation.Type, newLocation.Type)
}

func Test_FindAirport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock file
	file, err := ioutil.TempFile("", "airports")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(mockAirportJson)

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	airport, err := srv.FindAirport(dummyIataCode, file.Name())

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, dummyLocation.IataCode, airport.IataCode)
	assert.Equal(t, dummyLocation.Type, airport.Type)
}
