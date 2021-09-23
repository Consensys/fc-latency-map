package geomgr

//go:generate mockgen -destination mocks.go -package geomgr . GeoMgr

type GeoMgr interface {
	IPGeolocation(ip string) (float64, float64)
}
