package jobs

import (
	"fmt"
	"time"

	"github.com/ConsenSys/fc-latency-map/manager/constants"
	"github.com/ConsenSys/fc-latency-map/manager/export"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
)

func RunTaskGetMeasures() {
	locations.NewLocationHandler().UpdateLocations(constants.AirportTypeLarge)
	miners.NewMinerHandler().MinersParseStateMarket()
	probes.NewProbeHandler().Update()

	measrHdlr := measurements.NewHandler()
	measrHdlr.CreateMeasurements([]string{})
	measrHdlr.ImportMeasures()

	fn := fmt.Sprintf("data/exports/data_%v.json", time.Now().Unix())
	export.NewExportHandler().Export(fn)
}
