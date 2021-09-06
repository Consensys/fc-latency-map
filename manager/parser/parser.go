package parser

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
)

var mgrConfig = config.Config()
var nodeUrl string = mgrConfig.GetString("FILECOIN_NODE_URL")

func getMinersIP() {
	fMgr, err := filecoinmgr.NewFilecoinImpl(nodeUrl)
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
	fMgr.ExportJSON(minersWithIPs)
}


func Parse() error {
	getMinersIP()
}