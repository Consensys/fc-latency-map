package measurements

import (
	"fmt"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
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
		Handler: NewHandler(),
		Export:  export.NewExportHandler(),
	}
}

// completes the input
func (cmd *MesuresCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: measuresList, Description: "Get last measures"},
		{Text: measuresCreate, Description: "Create measurements"},
		{Text: measuresGet, Description: "Start getting measurements"},
		{Text: measuresExport, Description: "Export a json filename. ex: results_2021-09-17-17-17-00.json"},
	}
}

// executes the command
func (cmd *MesuresCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	case measuresCreate:
		cmd.Handler.CreateMeasurements()
	case measuresGet:
		cmd.Handler.ImportMeasures()
	case measuresList:
		cmd.measuresList(blocks)
	case measuresExport:
		cmd.measuresExport(blocks)
	default:
		log.Printf("unknown command: %s\n", blocks[0])
	}
}

func (cmd *MesuresCommander) measuresList(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing limit number")
		return
	}
}

func (cmd *MesuresCommander) measuresExport(blocks []string) {
	var fn string
	if len(blocks) == 1 {
		fn = fmt.Sprintf("data_%v.json", time.Now().Unix())
	}
	cmd.Export.Export(fn)
}
