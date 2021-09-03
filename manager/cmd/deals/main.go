package main

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
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
	fMgr.GetVerifiedDeals(blockHeight, 20)
}