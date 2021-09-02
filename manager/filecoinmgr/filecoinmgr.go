package filecoinmgr

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
)

type FilecoinMgr interface {

	// GetChainHead return chainhead data
	GetChainHead() (*types.TipSet, error)

	// GetChainHead return chainhead block height
	GetBlockHeight() (abi.ChainEpoch, error)

	// GetVerifiedDeals return verified deals for a range of block
	GetVerifiedDeals(height abi.ChainEpoch, offset uint) ([]VerifiedDeal, error)

}