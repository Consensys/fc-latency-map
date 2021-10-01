package seeds

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
)

const (
	seedData = "seed-data"
)

type SeederCommander struct {
}

func NewSeederCommander() cli.Commander {
	return &SeederCommander{}
}

// completes the input
func (cmd *SeederCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: seedData, Description: "Seed data(location, probes, miners, measures)"},
	}
}

// executes the command
func (cmd *SeederCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	case seedData:
		cmd.seedData()
	default:
		log.Printf("unknown command: %s\n", blocks[0])
	}
}

func (cmd *SeederCommander) seedData() {
	log.Println("Seed data ...")
	Seed()
}
