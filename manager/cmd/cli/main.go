package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	_ "gorm.io/driver/sqlite"

	"github.com/ConsenSys/fc-latency-map/manager/export"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
	"github.com/ConsenSys/fc-latency-map/manager/seeds"
)

const (
	locationsList   = "locations-list"
	locationsUpdate = "locations-update"
	locationsAdd    = "locations-add"
	locationsDelete = "locations-delete"
	probesUpdate    = "probes-update"
	probesList      = "probes-list"
	measuresGet     = "measures-get"
	measuresCreate  = "measures-create"
	measuresList    = "measures-list"
	measuresExport  = "measures-export"
	minersList      = "miners-list"
	minersUpdate    = "miners-update"
	minersParse     = "miners-parse"
	seedData        = "seed-data"
)

type LatencyMapCLI struct {
	probes       probes.ProbeHandler
	locations    locations.LocationHandler
	miners       miners.MinerHandler
	measurements measurements.Handler
	export       export.ExportHandler
}

// Start Client CLI
func main() {
	c := &LatencyMapCLI{
		probes:       *probes.NewProbeHandler(),
		locations:    *locations.NewLocationHandler(),
		miners:       *miners.NewMinerHandler(),
		measurements: *measurements.NewHandler(),
		export:       *export.NewExportHandler(),
	}

	if len(os.Args) > 1 {
		c.executor(strings.Join(os.Args[1:], " "))
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
		{Text: seedData, Description: "Seed data(location, probes, miners, measures)"},
		{Text: locationsList, Description: "List all locations"},
		{Text: locationsUpdate, Description: "insert airports in database. options: large / medium / small"},
		{Text: locationsAdd, Description: "Add location by country code. ex: location-add <country_code>"},
		{Text: locationsDelete, Description: "Delete location by country code. ex: location-delete <country_code>"},

		// probes
		{Text: probesUpdate, Description: "Update probes list by finding online and active probes"},
		{Text: probesList, Description: "Get probes list"},

		// measurements
		{Text: measuresList, Description: "Get last measures"},
		{Text: measuresCreate, Description: "Create measurements"},
		{Text: measuresGet, Description: "Start getting measurements"},
		{Text: measuresExport, Description: "Export a json filename. ex: results_2021-09-17-17-17-00.json"},

		// miners
		{Text: minersList, Description: "List all miners"},
		{Text: minersUpdate, Description: "Update miners list by finding active deals in past block heights. Offset is optional. ex: miners-update <offset>"},
		{Text: minersParse, Description: "Update miners list by finding active deals in a given block height. ex: miners-parse <block_height>"},

		// exit
		{Text: "exit", Description: "Exit the program"},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

//nolint:funlen
// executor executes the command
func (c *LatencyMapCLI) executor(in string) {
	log.Println("executor ", in)
	blocks := strings.Split(strings.TrimSpace(in), " ")

	log.Printf("Command: %s\n", blocks[0])

	switch blocks[0] {
	// Locations list
	case locationsList:
		c.locations.DisplayLocations()

	// Locations update
	case locationsUpdate:
		airportType := "large"
		if len(blocks) == 2 {
			airportType = blocks[1]
		}
		err := c.locations.UpdateLocations(airportType)
		if err != nil {
			log.Errorf("Error: %s\n", err)
		}

	// New location
	case locationsAdd:
		c.locationsAdd(blocks)

	// Delete location
	case locationsDelete:
		c.locationsDelete(blocks)

		// probes
	case probesUpdate:
		c.probes.Update()

	case probesList:
		c.probes.GetAllProbes()

		// Measurements
	case measuresCreate:
		c.measurements.CreateMeasurements()

	case measuresGet:
		c.measurements.ImportMeasures()

	case measuresList:
		c.measuresList(blocks)

	case measuresExport:
		c.measuresExport(blocks)

	case minersList:
		c.miners.GetAllMiners()

	case minersUpdate:
		c.minersUpdate(blocks)

	case minersParse:
		c.minersParse(blocks)

	case seedData:
		c.seedData()

	case "exit":
		log.Println("Shutdown ...")
		log.Println("Bye!")
		os.Exit(0)

	default:
		log.Printf("unknown command: %s\n", blocks[0])
	}
}

func (c *LatencyMapCLI) seedData() {
	log.Println("Seed data ...")
	seeds.Seed()
}

func (c *LatencyMapCLI) locationsAdd(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing location to add")
		return
	}

	log.Println("Add a location")
	location, err := c.locations.AddLocation(blocks[1])
	if err != nil {
		log.Error(err)
		return
	}
	log.Printf("ID: %d\n", location.ID)
}

func (c *LatencyMapCLI) locationsDelete(blocks []string) {
	if len(blocks) == 1 {
		log.Println("missing location to delete")
		return
	}

	log.Println("Delete a location")
	c.locations.DeleteLocation(blocks[1])
}

func (c *LatencyMapCLI) measuresList(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing limit number")
		return
	}
}

func (c *LatencyMapCLI) measuresExport(blocks []string) {
	var fn string
	if len(blocks) == 1 {
		fn = fmt.Sprintf("data_%v.json", time.Now().Unix())
	}
	c.export.Export(fn)
}

func (c *LatencyMapCLI) minersUpdate(blocks []string) {
	blockHeight := ""
	if len(blocks) > 1 {
		blockHeight = blocks[1]
	}
	c.miners.MinersUpdate(blockHeight)
}

func (c *LatencyMapCLI) minersParse(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing block height")

		return
	}
	height, err := strconv.ParseInt(blocks[1], 10, 64)
	if err != nil {
		log.Println("Error: provided block height is not a valid integer")

		return
	}
	c.miners.MinersParse(height)
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
