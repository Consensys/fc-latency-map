package parser

import (
	"fmt"
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/miners"
	"github.com/spf13/viper"
)

func Parse(
	conf *viper.Viper,
	fMgr filecoinmgr.FilecoinMgr,
	mSer miners.MinerService,
) error {
	getMinersIP(conf, fMgr, mSer)
	return nil
}

func getMinersIP(
	conf *viper.Viper,
	fMgr filecoinmgr.FilecoinMgr,
	mSer miners.MinerService,
) {
	blockHeight, err := fMgr.GetBlockHeight()
	if err != nil {
		log.Fatalf("get block failed: %s", err)
	}
	fmt.Printf("blockHeight: %+v\n", blockHeight)
	offset := conf.GetUint("FILECOIN_BLOCKS_OFFSET")
	verifiedDeals, err := fMgr.GetVerifiedDeals(blockHeight, offset)
	if err != nil {
		log.Fatalf("get block failed: %s", err)
	}
	minersWithIPs := mSer.GetMinerIPs(verifiedDeals)
	fmt.Printf("miners with IPs: %+v\n", minersWithIPs)
}
