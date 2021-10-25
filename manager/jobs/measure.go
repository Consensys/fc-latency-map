package jobs

import (
	"time"

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
	const ctx = "context"
	const crtMeas = "CreateMeasures"

	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Printf("Task started at %s", time.Now())

	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Println("Update locations...")
	err := locations.NewLocationHandler().UpdateLocations(constants.AirportTypeLarge)
	if err != nil {
		log.WithFields(log.Fields{
			ctx: crtMeas,
		}).Errorf("Error: %s", err)
	}

	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Println("Parse miners...")
	miners.BuildMinerHandlerInstance().MinersParseStateMarket()

	probeHdlr := probes.NewProbeHandler()
	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Println("Import probes...")
	probeHdlr.Import()
	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Println("Update probes...")
	probeHdlr.Update()

	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Println("Create measurements...")
	measurements.NewHandler().CreateMeasurements([]string{})

	log.WithFields(log.Fields{
		ctx: crtMeas,
	}).Printf("Task ended at %s", time.Now())
}

func RunTaskImportMeasures() {
	log.WithFields(log.Fields{
		"context": "ImportMeasures",
	}).Printf("Task started at %s", time.Now())

	conf := config.NewConfig()

	log.WithFields(log.Fields{
		"context": "ImportMeasures",
	}).Println("Import measurements...")
	measurements.NewHandler().ImportMeasures()

	log.WithFields(log.Fields{
		"context": "ImportMeasures",
	}).Println("Export data...")
	files := export.NewExportHandler().Export()

	log.WithFields(log.Fields{
		"context": "ImportMeasures",
	}).Println("Notify...")
	notif := webhook.NewNotifier(conf)
	notif.Notify(files)

	log.WithFields(log.Fields{
		"context": "ImportMeasures",
	}).Printf("Task ended at %s", time.Now())
}
