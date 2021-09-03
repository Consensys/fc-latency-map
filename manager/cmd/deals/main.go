package main

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/filecoinmgr"
)

func main() {
	fMgr, err := filecoinmgr.NewFilecoinImpl("https://node.glif.io/space07/lotus/rpc/v0")
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	blockHeight, err := fMgr.GetBlockHeight()
	if err != nil {
		log.Fatalf("get block failed: %s", err)
	}
	fmt.Printf("blockHeight: %+v\n", blockHeight)
	verifiedDeals, err := fMgr.GetVerifiedDeals(blockHeight, 100)
	if err != nil {
		log.Fatalf("get block failed: %s", err)
	}
	minersWithIPs := fMgr.GetMinerIPs(verifiedDeals)
	fmt.Printf("miners with IPs: %+v\n", minersWithIPs)
}
