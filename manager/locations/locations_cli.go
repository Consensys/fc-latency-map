package locations

import (
	"strings"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
	"github.com/ConsenSys/fc-latency-map/manager/constants"
)

const (
	locationsList   = "locations-list"
	locationsUpdate = "locations-update"
	locationsAdd    = "locations-add"
	locationsDelete = "locations-delete"
)

type LocationCommander struct {
	Handler *LocationHandler
}

func NewLocationCommander() cli.Commander {
	return &LocationCommander{
		Handler: NewLocationHandler(),
	}
}

// Complete completes the input
func (cmd *LocationCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: locationsList, Description: "List all locations"},
		{Text: locationsUpdate, Description: "insert airports in database. options: large / medium / small"},
		{Text: locationsAdd, Description: "Add location by country code. ex: location-add <country_code>"},
		{Text: locationsDelete, Description: "Delete location by country code. ex: location-delete <country_code>"},
	}
}

// Execute executes the command
func (cmd *LocationCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	// Locations list
	case locationsList:
		cmd.Handler.DisplayAllLocations()

	// Locations update
	case locationsUpdate:
		cmd.locationsUpdate(blocks)

	// New location
	case locationsAdd:
		cmd.locationsAdd(blocks)

	// Delete location
	case locationsDelete:
		cmd.locationsDelete(blocks)
	default:
		log.Printf("unknown command: %s\n", blocks[0])
	}
}

func (cmd *LocationCommander) locationsUpdate(blocks []string) {
	airportType := constants.AirportTypeLarge
	if len(blocks) == 2 {
		airportType = blocks[1]
	}
	err := cmd.Handler.UpdateLocations(airportType)
	if err != nil {
		log.Errorf("Error: %s\n", err)
	}
}

func (cmd *LocationCommander) locationsAdd(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing location to add")
		return
	}

	log.Println("Add a location")
	location, err := cmd.Handler.AddLocation(blocks[1])
	if err != nil {
		log.Error(err)
		return
	}
	log.Printf("ID: %d\n", location.ID)
}

func (cmd *LocationCommander) locationsDelete(blocks []string) {
	if len(blocks) == 1 {
		log.Println("missing location to delete")
		return
	}

	log.Println("Delete a location")
	cmd.Handler.DeleteLocation(blocks[1])
}
