package miners

import (
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
)

const (
	minersList             = "miners-list"
	minersParseOffset      = "miners-parse-offset"
	minersParseBlock       = "miners-parse-block"
	minersParseStateMarket = "miners-parse-state-market"
)

type MinerCommander struct {
	Handler *MinerHandler
}

func NewMinerCommander() cli.Commander {
	return &MinerCommander{
		Handler: NewMinerHandler(),
	}
}

// Complete completes the input
func (cmd *MinerCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: minersList, Description: "List all miners"},
		{Text: minersParseOffset, Description: "Parse miners by finding active deals in past block heights. Offset is optional. ex: miners-parse-offset <offset>"},
		{Text: minersParseBlock, Description: "Parse miners by finding active deals in a given block height. ex: miners-parse-block <block_height>"},
		{Text: minersParseStateMarket, Description: "Parse miners by finding state market deals. ex: miners-parse-state-market"},
	}
}

// Execute executes the command
func (cmd *MinerCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	case minersList:
		cmd.Handler.GetAllMiners()
	case minersParseOffset:
		cmd.minersParseOffset(blocks)
	case minersParseStateMarket:
		cmd.minersParseStateMarket()
	case minersParseBlock:
		cmd.minersParseBlock(blocks)
	default:
		log.Printf("unknown command: %s\n", blocks[0])
	}
}

func (cmd *MinerCommander) minersParseOffset(blocks []string) {
	blockHeight := ""
	if len(blocks) > 1 {
		blockHeight = blocks[1]
	}
	cmd.Handler.MinersParseOffset(blockHeight)
}

func (cmd *MinerCommander) minersParseBlock(blocks []string) {
	if len(blocks) == 1 {
		log.Println("Error: missing block height")

		return
	}
	height, err := strconv.ParseInt(blocks[1], 10, 64)
	if err != nil {
		log.Println("Error: provided block height is not a valid integer")

		return
	}
	cmd.Handler.MinersParseBlock(height)
}

func (cmd *MinerCommander) minersParseStateMarket() {
	cmd.Handler.minersParseStateMarket()
}
