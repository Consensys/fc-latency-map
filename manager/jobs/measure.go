package jobs

import (
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/constants"
	"github.com/ConsenSys/fc-latency-map/manager/export"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
	"github.com/ConsenSys/fc-latency-map/manager/webhook"
)

func RunTaskCreateMeasures() {
	log.Println("Update locations...")
	err := locations.NewLocationHandler().UpdateLocations(constants.AirportTypeLarge)
	if err != nil {
		log.Errorf("Error: %s\n", err)
	}

	log.Println("Parse miners...")
	miners.BuildMinerHandlerInstance().MinersParseStateMarket()

	probeHdlr := probes.NewProbeHandler()
	log.Println("Import probes...")
	probeHdlr.Import()
	log.Println("Update probes...")
	probeHdlr.Update()

	log.Println("Create measurements...")
	measurements.NewHandler().CreateMeasurements([]string{})
}

func RunTaskImportMeasures() {
	conf := config.NewConfig()

	log.Println("Import measurements...")
	measurements.NewHandler().ImportMeasures()

	log.Println("Export data...")
	files := export.NewExportHandler().Export()

	log.Println("Notify...")
	notif := webhook.NewNotifier(conf)
	notif.Notify(files)
}
