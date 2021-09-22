package locations

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationService interface {

	// GetLocations returns locations list
	DisplayLocations() []*models.Location

	// GetLocation returns a location
	GetLocation(location *models.Location) *models.Location

	// AddLocation creates a new location
	AddLocation(location *models.Location) *models.Location

	// DeleteLocation deletes a location
	DeleteLocation(location *models.Location) bool

	// CheckCountry checks if country exists
	CheckCountry(countryCode string) bool

	// FindAirport finds and returns airport
	FindAirport(airportCode string) (Airport, error)
}
