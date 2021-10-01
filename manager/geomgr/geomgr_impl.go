package geomgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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
		"ip":  ip,
		"geo": geo,
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
