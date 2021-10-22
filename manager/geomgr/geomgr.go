package geomgr

//go:generate mockgen -destination mocks.go -package geomgr . GeoMgr

type GeoMgr interface {
	IPGeolocation(ip string) (geolocation *Geolocation, err error)
}
