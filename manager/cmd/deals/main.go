package main

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/probes"
)

func main() {	

	bestProbes, err := probes.GetProbes()
	if err != nil {
		log.Fatalf("bestProbes failed: %s", err)
	}

	for i, probe := range bestProbes {

		fmt.Printf("probe nb: %v\n", i)
		fmt.Printf("ID: %v, Country: %v, IP: %v\n", probe.ID, probe.CountryCode, probe.AddressV4)
		
	}
}