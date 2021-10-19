package main

import (
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	_ "gorm.io/driver/sqlite"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
	"github.com/ConsenSys/fc-latency-map/manager/locations"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/ConsenSys/fc-latency-map/manager/probes"
)

type LatencyMapCLI struct {
	// Tip to maintain iteration order
	// See https://go.dev/blog/maps
	Commands   []string
	Commanders map[string]cli.Commander
}

// Start Client CLI
func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	c := &LatencyMapCLI{
		Commands: []string{"locations", "measures", "miners", "probes"},
		Commanders: map[string]cli.Commander{
			"locations": locations.NewLocationCommander(),
			"measures":  measurements.NewMesuresCommander(),
			"miners":    miners.NewMinerCommander(),
			"probes":    probes.NewProbeCommander(),
		},
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
		c.completer,
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}

// completer completes the input
func (c *LatencyMapCLI) completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	for _, command := range c.Commands {
		commander := c.Commanders[command]
		s = append(s, commander.Complete()...)
	}
	s = append(s, prompt.Suggest{Text: "exit", Description: "Exit the program"})
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *LatencyMapCLI) executor(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	log.Printf("Command: %s\n", blocks[0])

	if blocks[0] == "exit" {
		log.Println("Shutdown ...")
		log.Println("Bye!")
		os.Exit(0)
	}

	command := ""
	for _, candidate := range c.Commands {
		if strings.HasPrefix(blocks[0], candidate) {
			command = candidate
			break
		}
	}

	if commander, exists := c.Commanders[command]; exists {
		commander.Execute(in)
		log.Printf("Command ends: %s\n", blocks[0])
	} else {
		log.Printf("unknown command: %s\n", command)
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
