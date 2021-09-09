package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
)

const (
	locationList   = "location-list"
	locationAdd    = "location-add"
	locationDelete = "location-delete"
	probesUpdate   = "probes-update"
	measuresGet    = "measures-get"
	measuresList   = "measures-list"
	measuresExport = "measures-export"
	minersUpdate   = "miners-update"
)

type LatencyMapCLI struct {
	probes probes.Ripe
	locations locations.LocationServiceImpl
}

// Start Client CLI
func main() {
	probe, err := probes.NewClient("")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("starting probes")
	}
	c := &LatencyMapCLI{
		probes: *probe,
	}

	if len(os.Args) == 2 {
		c.executor(os.Args[1])
		os.Exit(0)
	}

	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("starting cli")
			debug.PrintStack()
		}

		handleExit()
	}()

	p := prompt.New(
		c.executor,
		completer,
		prompt.OptionPrefix(">>> "),
	)

	p.Run()
}

// completer completes the input
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		// location
		{Text: locationList, Description: "list all locations"},
		{Text: locationAdd, Description: "add location by country code. ex: location-add <country_code>"},
		{Text: locationDelete, Description: "delete location by country code. ex: location-delete <country_code>"},

		// probes
		{Text: probesUpdate, Description: "Update probes list by finding online and active probes"},

		// measurements
		{Text: measuresGet, Description: "start getting measurements"},
		{Text: measuresList, Description: "get last measures"},
		{Text: measuresExport, Description: "export a json filename. ex: results_2021-09-17-17-17-00.json"},

		// miners
		{Text: minersUpdate, Description: "Update miners list by find active deals in past blocks"},
		{Text: "exit", Description: "Exit the program"},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *LatencyMapCLI) executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	switch blocks[0] {
		
	// Location list
	case locationList:
		fmt.Printf("Command: %s \n", blocks[0])
		fmt.Println("List all location from db")
		locs := locations.NewLocationServiceImpl(dbMgr)
		locsList := locs.GetLocations()
		for _, location := range locsList {
			fmt.Printf("ID:%d - Country code: %s\n", location.ID, location.Country)
		}

	// New location
	case locationAdd:
		if len(blocks) == 1 {
			fmt.Println("missing location to add")
		}
		fmt.Printf("Command: %s \n", blocks[0])
		newLocation := models.Location{
			Country: blocks[1],
			Latitude:    "1.2",
			Longitude: "2.1",
		}
		locs := locations.NewLocationServiceImpl(dbMgr)
		newLocation = locs.AddLocation(newLocation)
		fmt.Printf("new location, ID:%d - Country code: %s\n", newLocation.ID, newLocation.Country)

	// Delete location
	case locationDelete:
		if len(blocks) == 1 {
			fmt.Println("missing location to delete")
		}
		fmt.Printf("Command: %s \n", blocks[0])
		location := models.Location{
			Country: blocks[1],
		}
		locs := locations.NewLocationServiceImpl(dbMgr)
		location = locs.GetLocation(location)
		if (location == models.Location{}) {
			fmt.Printf("Unable to find location %s\n", blocks[1])
		} else {
			locs.DeleteLocation(location)
			fmt.Printf("Location %d deleted\n", location.ID)
		}
		
		// probes
	case probesUpdate:
		fmt.Printf("Command: %s \n", blocks[0])
		c.probes.Update()

		// Measurements
	case measuresGet:
		fmt.Printf("Command: %s \n", blocks[0])

	case measuresList:
		if len(blocks) == 1 {
			fmt.Println("missing limit number")
		}
		fmt.Printf("Command: %s \n", blocks[0])

	case measuresExport:
		if len(blocks) == 1 {
			fmt.Println("missing filename")
		}
		fmt.Printf("Command: %s \n", blocks[0])
		fmt.Println("Get measures from db and export to a file")

	case minersUpdate:
		if len(blocks) == 1 {
			fmt.Println("add ")
		}
		fmt.Printf("Command: %s \n", blocks[0])
		fmt.Println("Call FC, get miners with active deals and store in db")

	case "exit":
		fmt.Println("Shutdown ...")
		fmt.Println("Bye!")
		os.Exit(0)

	default:
		fmt.Printf("unbknown command: %s\n", blocks[0])

	}
}

// handleExit fixes the problem of broken terminal when exit in Linux
// ref: https://www.gitmemory.com/issue/c-bata/go-prompt/228/820639887
func handleExit() {
	if _, err := os.Stat("/bin/stty"); os.IsNotExist(err) {
		return
	}
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	_ = rawModeOff.Wait()
}
