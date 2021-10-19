package geomgr

import (
	"fmt"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/stretchr/testify/assert"
)

var dummyIpAddress = "8.8.8.8"
var dummyLatitude = float64(37.6597400)
var dummyLongitude = float64(-97.5753300)
var dummyCountry = "US"

func Test_IPGeolocation_Fail_BadRequest(t *testing.T) {
	defer gock.Off()

	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)
	gock.New("http://www.geoplugin.net").
		Get("/json.gp").
		Reply(400)

	// Act
	lat, long, cntry := geo.IPGeolocation(dummyIpAddress)

	// Assert
	assert.NotNil(t, lat)
	assert.NotNil(t, long)
	assert.NotNil(t, cntry)
	assert.Equal(t, float64(0), lat)
	assert.Equal(t, float64(0), long)
	assert.Equal(t, "", cntry)
}

func Test_IPGeolocation_Fail_EmptyResponse(t *testing.T) {
	defer gock.Off()

	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)
	gock.New("http://www.geoplugin.net").
		Get("/json.gp").
		Reply(200)

	// Act
	lat, long, cntry := geo.IPGeolocation(dummyIpAddress)

	// Assert
	assert.NotNil(t, lat)
	assert.NotNil(t, long)
	assert.Equal(t, float64(0), lat)
	assert.Equal(t, float64(0), long)
	assert.Equal(t, "", cntry)
}

func Test_IPGeolocation_Fail_WrongJSON(t *testing.T) {
	defer gock.Off()

	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)
	gock.New("http://www.geoplugin.net").
		Get("/json.gp").
		Reply(200).
		JSON(map[string]interface{}{
			"status": 200,
		})

	// Act
	lat, long, cntry := geo.IPGeolocation(dummyIpAddress)

	// Assert
	assert.NotNil(t, lat)
	assert.NotNil(t, long)
	assert.Equal(t, float64(0), lat)
	assert.Equal(t, float64(0), long)
	assert.Equal(t, "", cntry)
}

func Test_IPGeolocation_OK(t *testing.T) {
	defer gock.Off()

	// Arrange
	mockConfig := config.NewMockConfig()
	geo := NewGeoMgrImpl(mockConfig)
	gock.New("http://www.geoplugin.net").
		Get("/json.gp").
		Reply(200).
		JSON(map[string]interface{}{
			"geoplugin_status":      200,
			"geoplugin_city":        "Unknown",
			"geoplugin_region":      "Kansas",
			"geoplugin_countryCode": dummyCountry,
			"geoplugin_latitude":    fmt.Sprintf("%f", dummyLatitude),
			"geoplugin_longitude":   fmt.Sprintf("%f", dummyLongitude),
			"geoplugin_timezone":    "America/Chicago",
		})

	// Act
	lat, long, cntry := geo.IPGeolocation(dummyIpAddress)

	// Assert
	assert.NotNil(t, lat)
	assert.NotNil(t, long)
	assert.Equal(t, dummyLatitude, lat)
	assert.Equal(t, dummyLongitude, long)
	assert.Equal(t, dummyCountry, cntry)
}
