package parser

import (
	"github.com/ConsenSys/fc-latency-map/manager/miners"
)

func Parse(mSer miners.MinerService) {
	mSer.ParseMiners(0)
}
