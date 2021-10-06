package probes

import (
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
	gomock "github.com/golang/mock/gomock"
)

var dummyLocation = models.Location{
	Name:   		"Charles de Gaulle International Airport",
	Country:   "FR",
	IataCode:  "CDG",
	Latitude:  49.012798,
	Longitude: 2.55,
	Type:      "large_airport",
}

var dummyProbeID = 42
var dummyCountryCode = "FR"
var dummyIataCode = "CDG"
var dummyLatitude = 49.012798
var dummyLongitude = 2.55


var dummyProbe = models.Probe{
	ProbeID:     dummyProbeID,
	CountryCode: dummyCountryCode,
	IataCode:    dummyIataCode,
	Location:    dummyLocation,
	Latitude:    dummyLatitude,
	Longitude:   dummyLongitude,
}



func Test_GetAllProbes_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()

	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	ripeMgr, err := ripemgr.NewRipeImpl(mockConfig)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr)

	// Act
	probes := srv.GetAllProbes()

	// Assert
	assert.NotNil(t, probes)
	assert.Empty(t, probes)
}

func Test_GetAllProbes_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()

	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	ripeMgr, err := ripemgr.NewRipeImpl(mockConfig)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr)

	// Act
	mockDbMgr.GetDB().Create(&([]*models.Location{&dummyLocation}))
	dummyProbe.Location = dummyLocation

	mockDbMgr.GetDB().Create(&([]*models.Probe{&dummyProbe}))
	probes := srv.GetAllProbes()

	// Assert
	assert.NotNil(t, probes)
	assert.NotEmpty(t, probes)

	actual := *(probes[0])
	assert.Equal(t, dummyProbe.ProbeID, actual.ProbeID)
	assert.Equal(t, dummyProbe.CountryCode, actual.CountryCode)
	assert.Equal(t, dummyProbe.IataCode, actual.IataCode)
	// assert.Equal(t, dummyProbe.Location, actual.Location)
	assert.Equal(t, dummyProbe.Latitude, actual.Latitude)
	assert.Equal(t, dummyProbe.Longitude, actual.Longitude)
}

func Test_GetTotalProbes_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()

	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	ripeMgr, err := ripemgr.NewRipeImpl(mockConfig)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr)

	// Act
	mockDbMgr.GetDB().Create(&([]*models.Probe{&dummyProbe}))
	count := srv.GetTotalProbes()

	// Assert
	assert.Equal(t, int64(1), count)
}