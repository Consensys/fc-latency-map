package filecoinmgr

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/filecoin-project/go-address"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/ipfs/go-cid"
	ma "github.com/multiformats/go-multiaddr"
)

type VerifiedDeal struct {
	MessageCid cid.Cid
	Provider   address.Address
}

type MinerIp struct {
	Provider address.Address
	MinerIPs []string
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

func (fMgr *FilecoinMgrImpl) GetVerifiedDeals(height abi.ChainEpoch, offset uint) ([]VerifiedDeal, error) {
	verifiedDeals := []VerifiedDeal{}
	for i := height - abi.ChainEpoch(offset); i <= height; i++ {
		fmt.Printf("Block number: %v (%v / %v)\n", i, height-i, offset)
		blockCids, _ := fMgr.api.ChainGetTipSetByHeight(context.Background(), abi.ChainEpoch(i), types.TipSetKey{})

		for _, cid := range blockCids.Cids() {
			messages, err := fMgr.api.ChainGetBlockMessages(context.Background(), cid)
			if err != nil {
				return []VerifiedDeal{}, err
			}
			for _, message := range messages.BlsMessages {
				// Method 4 is PublishStorageDeals
				if message.Method == 4 {
					var params market.PublishStorageDealsParams
					err = params.UnmarshalCBOR(bytes.NewReader(message.Params))
					if err != nil {
						return []VerifiedDeal{}, err
					}

					for _, deal := range params.Deals {
						proposal := deal.Proposal
						if proposal.VerifiedDeal {

							// TODO: Get deal Id
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
				}
			}
		}
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

func (fMgr *FilecoinMgrImpl) GetMinerIPs(verifiedDeals []VerifiedDeal) []MinerIp {
	var m = []MinerIp{}
	for _, deal := range verifiedDeals {
		provider := deal.Provider
		minerInfo, err := fMgr.api.StateMinerInfo(context.Background(), provider, types.TipSetKey{})
		if err != nil {
			continue
		}
		fmt.Printf("minerInfo: %+v\n", minerInfo)
		m = append(m, MinerIp{
			Provider: deal.Provider,
			MinerIPs: ipAddress(fMgr.multiAddrs(deal.Provider)),
		})
	}
	return m
}

func (fMgr *FilecoinMgrImpl) multiAddrs(addresss address.Address) []ma.Multiaddr {
	var m []ma.Multiaddr

	info, _ := fMgr.api.StateMinerInfo(context.Background(), addresss, types.TipSetKey{})

	for _, v := range info.Multiaddrs {
		fmt.Printf("info: %+v\n", info)
		if a, err := ma.NewMultiaddrBytes(v); err == nil {
			m = append(m, a)
			fmt.Printf("multiAddr: %+v\n", a)
		}
	}
	return m
}

func ipAddress(a []ma.Multiaddr) []string {
	var ips []string
	for _, v := range a {
		if ip, err := v.ValueForProtocol(ma.P_IP4); err == nil {
			ips = append(ips, ip)
		} else if ip, err := v.ValueForProtocol(ma.P_IP6); err == nil {
			ips = append(ips, ip)
		}
	}
	return ips
}

func (fMgr *FilecoinMgrImpl) ExportJSON(data []MinerIp) {
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("miners.json", file, 0644)
}