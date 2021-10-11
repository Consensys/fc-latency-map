package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/ConsenSys/fc-latency-map/manager/constants"
	"github.com/ConsenSys/fc-latency-map/manager/export"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
)

func RunTaskCreateMeasures() {
	log.Println("Update locations ...")
	locations.NewLocationHandler().UpdateLocations(constants.AirportTypeLarge)
	log.Println("Parse miners ...")
	miners.NewMinerHandler().MinersParseStateMarket()

	probeHdlr := probes.NewProbeHandler()
	log.Println("Import probes ...")
	probeHdlr.Import()
	log.Println("Update probes ...")
	probeHdlr.Update()

	log.Println("Create measurements ...")
	measurements.NewHandler().CreateMeasurements([]string{})
}

func RunTaskImportMeasures() {
	log.Println("Import measurements ...")
	measurements.NewHandler().ImportMeasures()

	log.Println("Export data ...")
	fn := fmt.Sprintf("data/exports/data_%v.json", time.Now().Unix())
	export.NewExportHandler().Export(fn)
}
