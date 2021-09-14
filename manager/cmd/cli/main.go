package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
)

const (
	locationList   = "location-list"
	locationAdd    = "location-add"
	locationDelete = "location-delete"
	probesUpdate   = "probes-update"
	probesList     = "probes-list"
	measuresGet    = "measures-get"
	measuresCreate = "measures-create"
	measuresList   = "measures-list"
	measuresExport = "measures-export"
	minersUpdate   = "miners-update"
	minersParse    = "miners-parse"
)

type LatencyMapCLI struct {
	probes       probes.ProbeHandler
	locations    locations.LocationHandler
	miners       miners.MinerHandler
	measurements measurements.Handler
}

// Start Client CLI
func main() {
	c := &LatencyMapCLI{
		probes:       *probes.NewProbeHandler(),
		locations:    *locations.NewLocationHandler(),
		miners:       *miners.NewMinerHandler(),
		measurements: *measurements.NewHandler(),
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
		{Text: locationList, Description: "List all locations"},
		{Text: locationAdd, Description: "Add location by country code. ex: location-add <country_code>"},
		{Text: locationDelete, Description: "Delete location by country code. ex: location-delete <country_code>"},

		// probes
		{Text: probesUpdate, Description: "Update probes list by finding online and active probes"},
		{Text: probesList, Description: "Get probes list"},

		// measurements
		{Text: measuresList, Description: "Get last measures"},
		{Text: measuresCreate, Description: "Create measurements"},
		{Text: measuresGet, Description: "Start getting measurements"},
		{Text: measuresExport, Description: "Export a json filename. ex: results_2021-09-17-17-17-00.json"},

		// miners
		{Text: minersUpdate, Description: "Update miners list by finding active deals in past block heights. Offset is optional. ex: miners-update <offset>"},
		{Text: minersParse, Description: "Update miners list by finding active deals in a given block height. ex: miners-parse <block_height>"},
		{Text: "exit", Description: "Exit the program"},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *LatencyMapCLI) executor(in string) {
	fmt.Println("executor ", in)
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {

	// Locations list
	case locationList:
		fmt.Printf("Command: %s \n", blocks[0])
		fmt.Println("List all location from db")
		c.locations.GetLocations()

	// New location
	case locationAdd:
		if len(blocks) == 1 {
			fmt.Println("Error: missing location to add")
			return
		}
		fmt.Printf("Command: %s \n", blocks[0])
		fmt.Println("Add a location")
		c.locations.AddLocation(blocks[1])

	// Delete location
	case locationDelete:
		if len(blocks) == 1 {
			fmt.Println("missing location to delete")
			return
		}
		fmt.Printf("Command: %s \n", blocks[0])
		fmt.Println("Delete a location")
		c.locations.DeleteLocation(blocks[1])

		// probes
	case probesUpdate:
		fmt.Printf("Command: %s \n", blocks[0])
		c.probes.Update()

	case probesList:
		fmt.Printf("Command: %s \n", blocks[0])
		c.probes.GetAllProbes()

		// Measurements
	case measuresCreate:
		c.measurements.CreateMeasurements()

	case measuresGet:
		fmt.Printf("Command: %s \n", blocks[0])

		c.measurements.GetMeasures()

	case measuresList:
		if len(blocks) == 1 {
			fmt.Println("Error: missing limit number")
			return
		}
		fmt.Printf("Command: %s \n", blocks[0])

	case measuresExport:
		if len(blocks) == 1 {
			fmt.Println("Error: missing filename")
			return
		}
		c.measurements.ExportData(blocks[1])

	case minersUpdate:
		fmt.Printf("Command: %s \n", blocks[0])
		blockHeight := ""
		if len(blocks) > 1 {
			blockHeight = blocks[1]
		}
		c.miners.MinersUpdate(blockHeight)

	case minersParse:
		fmt.Printf("Command: %s \n", blocks[0])
		if len(blocks) == 1 {
			fmt.Println("Error: missing block height")
			return
		}
		height, err := strconv.ParseInt(blocks[1], 10, 64)
		if err != nil {
			fmt.Println("Error: provided block height is not a valid integer")
			return
		}
		c.miners.MinersParse(height)

	case "exit":
		fmt.Println("Shutdown ...")
		fmt.Println("Bye!")
		os.Exit(0)

	default:
		fmt.Printf("unknown command: %s\n", blocks[0])

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
