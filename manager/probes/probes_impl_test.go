package probes

import (
	"errors"
	"log"
	"testing"

	gomock "github.com/golang/mock/gomock"
	atlas "github.com/keltia/ripe-atlas"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/geomgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

var dummyProbeID = 42
var dummyCountryCode = "FR"
var dummyLatitude = 49.012798
var dummyLongitude = 2.55

var dummyProbe = models.Probe{
	ProbeID:     dummyProbeID,
	CountryCode: dummyCountryCode,
	Latitude:    dummyLatitude,
	Longitude:   dummyLongitude,
}

var mockOpts = map[string]string{
	"status_name": "Connected",
	"is_public":   "true",
	"sort":        "id",
}

func Test_ListProbes_Empty(t *testing.T) {
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
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr, nil)

	// Act
	probes := srv.ListProbes()

	// Assert
	assert.Equal(t, 0, len(probes))
}

func Test_ListProbes_OK(t *testing.T) {
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
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr, nil)

	// Act
	mockDbMgr.GetDB().Create(&([]*models.Probe{&dummyProbe}))
	probes := srv.ListProbes()

	// Assert
	assert.Equal(t, 1, len(probes))
}

func Test_GetTotalProbes_Empty(t *testing.T) {
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
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr, nil)

	// Act
	count := srv.GetTotalProbes()

	// Assert
	assert.Equal(t, int64(0), count)
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
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr, nil)

	// Act
	mockDbMgr.GetDB().Create(&([]*models.Probe{&dummyProbe}))
	count := srv.GetTotalProbes()

	// Assert
	assert.Equal(t, int64(1), count)
}

func Test_RequestProbes_OK(t *testing.T) {
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
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr, nil)

	// Act
	errr := srv.RequestProbes()

	// Assert
	assert.Nil(t, errr)
}

func Test_Update_OK(t *testing.T) {
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
	srv, _ := NewProbeServiceImpl(mockDbMgr, ripeMgr, nil)

	// Act
	updated := srv.Update()

	// Assert
	assert.True(t, updated)
}

func Test_ImportProbes_Fail_ProbesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockDbMgr := db.NewMockDatabaseMgr()
	mockRipeMgr := ripemgr.NewMockRipeMgr(ctrl)
	mockGeoMgr := geomgr.NewMockGeoMgr(ctrl)
	srv, _ := NewProbeServiceImpl(mockDbMgr, mockRipeMgr, mockGeoMgr)
	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	// Act
	mockRipeMgr.EXPECT().GetProbes(mockOpts).Return(nil, errors.New(""))
	imported := srv.ImportProbes()

	// Assert
	assert.False(t, imported)
}

func Test_ImportProbes_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockDbMgr := db.NewMockDatabaseMgr()
	mockRipeMgr := ripemgr.NewMockRipeMgr(ctrl)
	mockGeoMgr := geomgr.NewMockGeoMgr(ctrl)
	mockProbes := make([]atlas.Probe, 0)
	srv, _ := NewProbeServiceImpl(mockDbMgr, mockRipeMgr, mockGeoMgr)
	sqlDB, _ := mockDbMgr.GetDB().DB()
	defer sqlDB.Close()

	// Act
	mockRipeMgr.EXPECT().GetProbes(mockOpts).Return(mockProbes, nil)
	imported := srv.ImportProbes()

	// Assert
	assert.True(t, imported)
}
