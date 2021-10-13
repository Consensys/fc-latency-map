package geomgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GeoMgrImpl struct {
	conf *viper.Viper
}

const sleepTime = 510 * time.Millisecond

func NewGeoMgrImpl(v *viper.Viper) GeoMgr {
	return &GeoMgrImpl{
		conf: v,
	}
}

func (g *GeoMgrImpl) IPGeolocation(ip string) (lat, long float64) {
	// to free use of api
	// you can check you ip status https://www.geoplugin.com/ip_status.php?ip=xx.xx.xx.xx

	time.Sleep(sleepTime)
	response, err := http.Get(g.ipgeoURL(ip)) //nolint:noctx
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ipgeolocation")

		return 0, 0
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ipgeolocation")

		return 0, 0
	}

	var geo GeoIP
	err = json.Unmarshal(body, &geo)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ipgeolocation")

		return 0, 0
	}

	log.WithFields(log.Fields{
		"ip":        ip,
		"region":    geo.GeopluginRegion,
		"city":      geo.GeopluginCity,
		"latitude":  geo.GeopluginLatitude,
		"longitude": geo.GeopluginLongitude,
		"status":    geo.GeopluginStatus,
		"timezone":  geo.GeopluginTimezone,
	}).Info("ipgeolocation")

	return toFloat(geo.GeopluginLatitude), toFloat(geo.GeopluginLongitude)
}

func (g *GeoMgrImpl) ipgeoURL(ip string) string {
	return fmt.Sprintf("http://www.geoplugin.net/json.gp?ip=%s", ip)
}

func toFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return f
	}
	return 0
}

func (g *GeoMgrImpl) FindCountry(lat, long float64) string {
	time.Sleep(sleepTime)
	response, err := http.Get(g.ipgeoLocationURL(lat, long)) //nolint:noctx
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ipgeolocation country")

		return ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ipgeolocation country")

		return ""
	}

	var l Location
	err = json.Unmarshal(body, &l)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ipgeolocation country")

		return ""
	}

	log.WithFields(log.Fields{
		"long":        lat,
		"lat":         long,
		"countryCode": l.GeopluginCountryCode,
	}).Info("ipgeolocation country")

	return l.GeopluginCountryCode
}

func (g *GeoMgrImpl) ipgeoLocationURL(lat, long float64) string {
	return fmt.Sprintf("http://www.geoplugin.net/extras/location.gp?lat=%v&long=%v&format=json", lat, long)
}
