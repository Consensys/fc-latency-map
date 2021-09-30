package locations

//go:generate mockgen -destination mocks.go -package locations . LocationService

import (
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationService interface {

	// GetAllLocations returns locations list
	GetAllLocations() []*models.Location

	// GetLocation returns a location
	GetLocation(location *models.Location) *models.Location

	// AddLocation creates a new location
	AddLocation(location *models.Location) *models.Location

	// DeleteLocation deletes a location
	DeleteLocation(location *models.Location) bool

	// CheckCountry checks if country exists
	CheckCountry(countryCode string) bool

	// UpdateLocations create airports in database
	UpdateLocations(airportType string, filename string) error

	// ExtractAirports returns airports
	ExtractAirports(filename string) ([]Airport, error)

	// FindAirport finds and returns airport
	FindAirport(airportCode string, filename string) (Airport, error)
}
