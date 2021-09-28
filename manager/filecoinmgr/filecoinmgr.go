package filecoinmgr

//go:generatez mockgen -destination mocks.go -package filecoinmgr . FilecoinMgr

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
)

type FilecoinMgr interface {

	// GetChainHead return chainhead data
	GetChainHead() (*types.TipSet, error)

	// GetBlockHeight return chainhead block height
	GetBlockHeight() (abi.ChainEpoch, error)

	// GetVerifiedDealsByBlockRange return verified deals for a range of block
	GetVerifiedDealsByBlockRange(height abi.ChainEpoch, offset int) ([]VerifiedDeal, error)

	// GetVerifiedDealsByBlockHeight return verified deals for a block height
	GetVerifiedDealsByBlockHeight(height abi.ChainEpoch) ([]VerifiedDeal, error)

	// GetMinerInfo
	GetMinerInfo(addr address.Address) (miner.MinerInfo, error)
}
