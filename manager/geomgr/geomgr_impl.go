package geomgr

import (
	"net"

	"github.com/oschwald/maxminddb-golang"

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

func (g *GeoMgrImpl) IPGeolocation(ip string) (geolocation *Geolocation, err error) {
	db, err := maxminddb.Open(g.conf.GetString("GEOLITE2_MMDB"))
	if err != nil {
		return &Geolocation{}, err
	}
	defer db.Close()

	var l IPGeolocation
	_, found, err := db.LookupNetwork(net.ParseIP(ip), &l)
	if err != nil {
		return &Geolocation{}, err
	}
	if !found {
		return &Geolocation{}, nil
	}

	return &Geolocation{
		Country:   l.Country.ISOCode,
		Latitude:  l.Location.Latitude,
		Longitude: l.Location.Longitude,
	}, nil
}
