package geomgr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GeoMgrImpl struct {
	conf *viper.Viper
}

func NewGeoMgrImpl(v *viper.Viper) GeoMgr {
	return &GeoMgrImpl{
		conf: v,
	}
}

func (g *GeoMgrImpl) IPGeolocation(ip string) (lat, long float64) {
	response, err := http.NewRequestWithContext(context.Background(), "GET", g.ipgeoURL(ip), nil)
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
		"ip":  ip,
		"geo": geo,
	}).Info("ipgeolocation")

	return geo.Latitude, geo.Longitude
}

func (g *GeoMgrImpl) ipgeoURL(ip string) string {
	key := g.conf.Get("IPGEOLOCATION_ABSTRACTAPI_KEY")

	return fmt.Sprintf("https://ipgeolocation.abstractapi.com/v1/?"+
		"&fields=city,city_geoname_id,country_code,continent,latitude,longitude"+
		"&api_key=%s"+
		"&ip_address=%s", key, ip)
}
