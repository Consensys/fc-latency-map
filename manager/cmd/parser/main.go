package main

import (
	"github.com/ConsenSys/fc-latency-map/manager/miners"
)

func main() {
	mHdl := miners.NewMinerHandler()
	mHdl.MinersUpdate("")
}
