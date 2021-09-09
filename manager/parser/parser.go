package parser

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
)

var mgrConfig = config.NewConfig()
var offset = mgrConfig.GetUint("FILECOIN_BLOCKS_OFFSET")

func Parse(
	fMgr filecoinmgr.FilecoinMgr,
	mSer miners.MinerService,
) error {
	getMinersIP(fMgr, mSer)
	return nil
}

func getMinersIP(
	fMgr filecoinmgr.FilecoinMgr,
	mSer miners.MinerService,
) {
	blockHeight, err := fMgr.GetBlockHeight()
	if err != nil {
		log.Fatalf("get block failed: %s", err)
	}
	fmt.Printf("blockHeight: %+v\n", blockHeight)
	verifiedDeals, err := fMgr.GetVerifiedDeals(blockHeight, offset)
	if err != nil {
		log.Fatalf("get block failed: %s", err)
	}
	minersWithIPs := mSer.GetMinerIPs(verifiedDeals)
	fmt.Printf("miners with IPs: %+v\n", minersWithIPs)
}
