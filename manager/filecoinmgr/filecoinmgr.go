package filecoinmgr

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
)

type FilecoinMgr interface {

	// GetChainHead return chainhead data
	GetChainHead() (*types.TipSet, error)

	// GetChainHead return chainhead block height
	GetBlockHeight() (abi.ChainEpoch, error)

	// GetVerifiedDeals return verified deals for a range of block
	GetVerifiedDeals(height abi.ChainEpoch, offset uint) ([]VerifiedDeal, error)

	// GetVerifiedDeals return verified deals for a block height
	GetVerifiedDealsByBlockHeight(height abi.ChainEpoch) ([]VerifiedDeal, error)

	// GetMinerInfo
	GetMinerInfo(addr address.Address) (miner.MinerInfo, error)
}
