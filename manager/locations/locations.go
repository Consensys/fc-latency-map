package locations

//go:generate mockgen -destination mocks.go -package locations . LocationService

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationService interface {

	// DisplayLocations returns locations list
	DisplayLocations() []*models.Location

	// UpdateLocations create airports in database
	UpdateLocations(airportType string) (bool, error)

	// GetLocation returns a location
	GetLocation(location *models.Location) *models.Location

	// AddLocation creates a new location
	AddLocation(location *models.Location) *models.Location

	// DeleteLocation deletes a location
	DeleteLocation(location *models.Location) bool

	// CheckCountry checks if country exists
	CheckCountry(countryCode string) bool

	// ExtractAirports returns airports
	ExtractAirports() ([]Airport, error)

	// FindAirport finds and returns airport
	FindAirport(airportCode string) (Airport, error)
}
