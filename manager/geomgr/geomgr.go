package geomgr

type GeoMgr interface {
	IPGeolocation(ip string) (float64, float64)
}
