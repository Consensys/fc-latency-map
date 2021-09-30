package geolocation

import (
	"math"

	"gonum.org/v1/gonum/spatial/vptree"
)

// Place is a vptree.Comparable implementations.
type Place struct {
	ID        int
	Latitude  float64
	Longitude float64
}

// Distance returns the distance between the receiver and c.
func (p Place) Distance(c vptree.Comparable) float64 {
	q := c.(Place)
	return Haversine(p.Latitude, p.Longitude, q.Latitude, q.Longitude)
}

// Haversine returns the distance between two geographic coordinates.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371 // km
	sdLat := math.Sin(radians(lat2-lat1) / 2)
	sdLon := math.Sin(radians(lon2-lon1) / 2)
	a := sdLat*sdLat + math.Cos(radians(lat1))*math.Cos(radians(lat2))*sdLon*sdLon
	d := 2 * r * math.Asin(math.Sqrt(a))
	return d // km
}

// radians convert degrees into radians
func radians(degree float64) float64 {
	const degrees180 = 180
	return degree * math.Pi / degrees180
}
