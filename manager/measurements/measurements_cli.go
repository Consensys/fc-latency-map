package measurements

import (
	"strings"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
	"github.com/ConsenSys/fc-latency-map/manager/export"
)

const (
	measuresGet    = "measures-get"
	measuresCreate = "measures-create"
	measuresList   = "measures-list"
	measuresExport = "measures-export"
)

type MesuresCommander struct {
	Handler *Handler
	Export  *export.ExportHandler
}

func NewMesuresCommander() cli.Commander {
	return &MesuresCommander{
		Handler: newHandler(),
		Export:  export.NewExportHandler(),
	}
}

// Complete completes the input
func (cmd *MesuresCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: measuresList, Description: "Get last measures"},
		{Text: measuresCreate, Description: "Create measurements"},
		{Text: measuresGet, Description: "Start getting measurements"},
		{Text: measuresExport, Description: "Export a json filename. ex: results_2021-09-17-17-17-00.json"},
	}
}

// Execute executes the command
func (cmd *MesuresCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	case measuresCreate:
		cmd.Handler.CreateMeasurements(blocks)
	case measuresGet:
		cmd.Handler.ImportMeasures()
	case measuresList:
		cmd.measuresList(blocks)
	case measuresExport:
		cmd.measuresExport()
	default:
		log.Printf("unknown command: %stopped\n", blocks[0])
	}
}

func (cmd *MesuresCommander) measuresList(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing limit number")
		return
	}
}

func (cmd *MesuresCommander) measuresExport() {
	cmd.Export.Export()
}
