package probes

import (
	"strings"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
)

const (
	probesUpdate = "probes-update"
	probesImport = "probes-import"
	probesList   = "probes-list"
)

type ProbeCommander struct {
	Handler *ProbeHandler
}

func NewProbeCommander() cli.Commander {
	return &ProbeCommander{
		Handler: NewProbeHandler(),
	}
}

// Complete completes the input
func (cmd *ProbeCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: probesUpdate, Description: "Update probes list by finding online and active probes"},
		{Text: probesList, Description: "Get probes list"},
	}
}

// Execute executes the command
func (cmd *ProbeCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	case probesUpdate:
		cmd.Handler.Update()
	case probesList:
		cmd.Handler.List()
	case probesImport:
		cmd.Handler.Import()
	default:
		log.Printf("unknown command: %s", blocks[0])
	}
}
