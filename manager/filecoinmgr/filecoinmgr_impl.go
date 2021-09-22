package filecoinmgr

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/ipfs/go-cid"
)

type VerifiedDeal struct {
	MessageCid cid.Cid
	Provider   address.Address
}

type FilecoinMgrImpl struct {
	addr string
	api  lotusapi.FullNodeStruct
}

func NewFilecoinImpl(addr string) (FilecoinMgr, error) {
	// Node connection
	headers := http.Header{}
	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		addr,
		"Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal},
		headers,
	)
	if err != nil {
		return nil, err
	}
	defer closer()

	return &FilecoinMgrImpl{
		addr: addr,
		api:  api,
	}, nil
}

func (fMgr *FilecoinMgrImpl) GetChainHead() (*types.TipSet, error) {
	tipset, err := fMgr.api.ChainHead(context.Background())
	if err != nil {
		return nil, err
	}
	return tipset, nil
}

func (fMgr *FilecoinMgrImpl) GetBlockHeight() (abi.ChainEpoch, error) {
	tipset, err := fMgr.GetChainHead()
	if err != nil {
		return 0, err
	}
	return tipset.Height(), err
}

func (fMgr *FilecoinMgrImpl) GetMinerInfo(addr address.Address) (miner.MinerInfo, error) {
	minerInfo, err := fMgr.api.StateMinerInfo(context.Background(), addr, types.TipSetKey{})
	if err != nil {
		return miner.MinerInfo{}, err
	}
	return minerInfo, err
}

func (fMgr *FilecoinMgrImpl) GetVerifiedDealsByBlockRange(height abi.ChainEpoch, offset uint) ([]VerifiedDeal, error) {
	verifiedDeals := []VerifiedDeal{}
	for i := height - abi.ChainEpoch(offset); i <= height; i++ {
		fmt.Printf("Block number: %v (%v / %v)\n", i, height-i, offset)

		newVerifiedDeals, err := fMgr.getVerifiedDealsByBlock(height)
		if err != nil {
			return []VerifiedDeal{}, err
		}
		verifiedDeals = append(verifiedDeals, newVerifiedDeals...)
	}

	fmt.Printf("verifiedDeals: %+v\n", verifiedDeals)
	return verifiedDeals, nil
}

func CheckIsVerifiedDeal(verifiedDeal VerifiedDeal, verifiedDeals []VerifiedDeal) bool {
	for _, deal := range verifiedDeals {
		if deal.MessageCid == verifiedDeal.MessageCid && deal.Provider == verifiedDeal.Provider {
			return true
		}
	}
	return false
}

func (fMgr *FilecoinMgrImpl) GetVerifiedDealsByBlockHeight(height abi.ChainEpoch) ([]VerifiedDeal, error) {
	fmt.Printf("Block number: %v\n", height)
	verifiedDeals, err := fMgr.getVerifiedDealsByBlock(height)
	if err != nil {
		return []VerifiedDeal{}, err
	}
	return verifiedDeals, nil
}

func (fMgr *FilecoinMgrImpl) getVerifiedDealsByBlock(height abi.ChainEpoch) ([]VerifiedDeal, error) {
	blockCids, _ := fMgr.api.ChainGetTipSetByHeight(context.Background(), height, types.TipSetKey{})
	verifiedDeals := []VerifiedDeal{}
	for _, cid := range blockCids.Cids() {
		messages, err := fMgr.api.ChainGetBlockMessages(context.Background(), cid)
		if err != nil {
			return []VerifiedDeal{}, err
		}
		for _, message := range messages.BlsMessages {
			// Method 4 is PublishStorageDeals
			if message.Method != 4 {
				continue
			}
			params := &market.PublishStorageDealsParams{}
			err = params.UnmarshalCBOR(bytes.NewReader(message.Params))
			if err != nil {
				return []VerifiedDeal{}, err
			}

			verifiedDeals = fMgr.getVerifiedDeals(params, message, verifiedDeals)
		}
	}

	fmt.Printf("verifiedDeals: %+v\n", verifiedDeals)
	return verifiedDeals, nil
}

func (fMgr *FilecoinMgrImpl) getVerifiedDeals(params *market.PublishStorageDealsParams, message *types.Message, verifiedDeals []VerifiedDeal) []VerifiedDeal {
	for _, deal := range params.Deals { //nolint:gocritic
		proposal := deal.Proposal

		if proposal.VerifiedDeal {
			verifiedDeal := VerifiedDeal{
				MessageCid: message.Cid(),
				Provider:   proposal.Provider,
			}
			if !CheckIsVerifiedDeal(verifiedDeal, verifiedDeals) {
				fmt.Println("Verified deal found")
				verifiedDeals = append(verifiedDeals, verifiedDeal)
			}
		}
	}
	return verifiedDeals
}
