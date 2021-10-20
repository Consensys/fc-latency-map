package geomgr

import (
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/stretchr/testify/assert"
)

var dummyIpAddress = "8.8.8.8"
var dummyPrivateIpAddress = "192.168.89.8"
var dummyLatitude = float64(37.6597400)
var dummyLongitude = float64(-97.5753300)
var dummyCountry = "US"

func Test_IPGeolocation_OK(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)

	// Act
	g, err := geo.IPGeolocation(dummyIpAddress)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, g.Latitude)
	assert.NotNil(t, g.Longitude)
	assert.NotNil(t, g.Longitude)
	assert.Equal(t, (37.751), g.Latitude)
	assert.Equal(t, (-97.822), g.Longitude)
	assert.Equal(t, "US", g.Country)
}

func Test_IPGeolocation_PrivateIP(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)

	// Act
	g, err := geo.IPGeolocation(dummyPrivateIpAddress)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, g.Latitude)
	assert.NotNil(t, g.Longitude)
	assert.NotNil(t, g.Longitude)
	assert.Equal(t, 0.0, g.Latitude)
	assert.Equal(t, 0.0, g.Longitude)
	assert.Equal(t, "", g.Country)
}

func Test_IPGeolocation_InvalidIP(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)

	// Act
	g, err := geo.IPGeolocation("dummyPrivateIpAddress")

	// Assert
	assert.NotNil(t, err)
	assert.NotNil(t, g.Latitude)
	assert.NotNil(t, g.Longitude)
	assert.NotNil(t, g.Longitude)
	assert.Equal(t, 0.0, g.Latitude)
	assert.Equal(t, 0.0, g.Longitude)
	assert.Equal(t, "", g.Country)
}

func Test_IPGeolocation_InvalidFile(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	mockConfig.Set("GEOLITE2_MMDB", "")
	geo := NewGeoMgrImpl(mockConfig)

	// Act
	g, err := geo.IPGeolocation(dummyIpAddress)

	// Assert
	assert.NotNil(t, err)
	assert.NotNil(t, g.Latitude)
	assert.NotNil(t, g.Longitude)
	assert.NotNil(t, g.Longitude)
	assert.Equal(t, 0.0, g.Latitude)
	assert.Equal(t, 0.0, g.Longitude)
	assert.Equal(t, "", g.Country)
}