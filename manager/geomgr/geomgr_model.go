package geomgr

type Geolocation struct {
	Country   string
	Latitude  float64
	Longitude float64
}

type Country struct {
	ISOCode string `maxminddb:"iso_code"`
}

type Location struct {
	Latitude float64 `maxminddb:"latitude"`
	// Longitude is directly nested within the parent map.
	Longitude float64 `maxminddb:"longitude"`
	// TimeZone is indirected via a pointer.
	TimeZoneOffset uintptr `maxminddb:"time_zone"`
}

type IPGeolocation struct {
	Country  Country  `maxminddb:"country"`
	Location Location `maxminddb:"location"`
}
