package locations

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

var dummyCountryError = "FRX"

var dummyCountry = "FR"
var dummyIataCode = "CDG"
var dummyGeoLatitude = 2.55
var dummyGeoLongitude = 49.012798
var dummyType = "large_airport"

var dummyBlockHeight = int64(42)
var dummyLocation = models.Location{
	Country:   	dummyCountry,
	IataCode:   dummyIataCode,
	Latitude:  	dummyGeoLatitude,
	Longitude: 	dummyGeoLongitude,
	Type: 			dummyType,
}

type TestSuite struct {
	suite.Suite
}

// before each test
func (suite *TestSuite) SetupTest() {
	mockDbMgr := db.NewMockDatabaseMgr()
	// mockDbMgr.GetDB().Delete(&models.Location{})
	mockDbMgr.GetDB().Where("1 = 1").Delete(&models.Location{})
}

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
	assert.Equal(t, dummyLocation.Country, actual.Country)
	assert.Equal(t, dummyLocation.IataCode, actual.IataCode)
	assert.Equal(t, dummyLocation.Latitude, actual.Latitude)
	assert.Equal(t, dummyLocation.Longitude, actual.Longitude)
	assert.Equal(t, dummyLocation.Type, actual.Type)

	
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

	fmt.Printf("===>>> %v\n\n", location)

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

func Test_CheckCountry_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	isCountry := srv.CheckCountry(dummyCountry)

	// Assert
	assert.Equal(t, true, isCountry)
}

func Test_CheckCountry_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	srv := NewLocationServiceImpl(mockConfig, mockDbMgr)

	// Act
	isCountry := srv.CheckCountry(dummyCountryError)

	// Assert
	assert.Equal(t, false, isCountry)
}
