package main

import (
    "fmt"

    "os"
    "os/exec"
    "runtime/debug"
    "strings"

    "github.com/c-bata/go-prompt"
    log "github.com/sirupsen/logrus"
)

type LatencyMapCLI struct {
}

// Start Client CLI
func main() {
    c := &LatencyMapCLI{}

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
        completer,
        prompt.OptionPrefix(">>> "),
    )

    p.Run()
}

// completer completes the input
func completer(d prompt.Document) []prompt.Suggest {
    s := []prompt.Suggest{
        // location
        {Text: "location-list", Description: "list all locations"},
        {Text: "location-add", Description: "add location by country code. ex: location-add <country_code>"},
        {Text: "location-delete", Description: "delete location by country code. ex: location-delete <country_code>"},

        // probes
        {Text: "probes-update", Description: "Update probes list by finding online and active probes"},

        // measurements
        {Text: "measures-get", Description: "start getting measurements"},
        {Text: "measures-list", Description: "get last measures"},
        {Text: "measures-export", Description: "export a json filename. ex: results_2021-09-17-17-17-00.json"},

        // miners
        {Text: "miners-update", Description: "Update miners list by find active deals in past blocks"},
        {Text: "exit", Description: "Exit the program"},
    }

    return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *LatencyMapCLI) executor(in string) {
    in = strings.TrimSpace(in)
    blocks := strings.Split(in, " ")

    switch blocks[0] {
    case "location-list":
        fmt.Printf("Command: %s \n", blocks[0])
        fmt.Println("List all location from db")

    case "location-add":
        if len(blocks) == 1 {
            fmt.Println("missing location to add")
        }
        fmt.Printf("Command: %s \n", blocks[0])

    case "location-delete":
        if len(blocks) == 1 {
            fmt.Println("missing location to delete")
        }
        fmt.Printf("Command: %s \n", blocks[0])
        // probes
    case "probes-update":
        fmt.Printf("Command: %s \n", blocks[0])

        // Measurements
    case "measures-get":
        fmt.Printf("Command: %s \n", blocks[0])

    case "measures-list":
        if len(blocks) == 1 {
            fmt.Println("missing limit number")
        }
        fmt.Printf("Command: %s \n", blocks[0])

    case "measures-export":
        if len(blocks) == 1 {
            fmt.Println("missing filename")
        }
        fmt.Printf("Command: %s \n", blocks[0])
        fmt.Println("Get measures from db and export to a file")

    case "miners-update":
        if len(blocks) == 1 {
            fmt.Println("add ")
        }
        fmt.Printf("Command: %s \n", blocks[0])
        fmt.Println("Call FC, get miners with active deals and store in db")

    case "exit":
        fmt.Println("Shutdown ...")
        fmt.Println("Bye!")
        os.Exit(0)

    default:
        fmt.Printf("unbknown command: %s\n", blocks[0])

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