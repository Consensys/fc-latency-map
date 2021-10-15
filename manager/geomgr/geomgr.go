package geomgr

//go:generate mockgen -destination mocks.go -package geomgr . GeoMgr

type GeoMgr interface {
	IPGeolocation(ip string) (lat float64, long float64, code string)
	FindCountry(lat, long float64) string
}
