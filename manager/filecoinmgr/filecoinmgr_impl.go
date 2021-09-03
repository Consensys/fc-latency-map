package filecoinmgr

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-address"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/ipfs/go-cid"
)

type VerifiedDeal struct {
	MessageCid cid.Cid
	Provider address.Address
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
		addr:         addr,
		api:         api,
		
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
		fmt.Printf("Block number: %v (%v / %v)\n", i, height - i, offset)
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
								Provider: proposal.Provider,
							}
							if (!CheckIsVerifiedDeal(verifiedDeal, verifiedDeals)) {
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
		if (deal.MessageCid == verifiedDeal.MessageCid && deal.Provider == verifiedDeal.Provider) {
			return true
		}
	}
	return false
}


// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/filecoin-project/go-address"
// 	jsonrpc "github.com/filecoin-project/go-jsonrpc"
// 	"github.com/filecoin-project/go-state-types/abi"
// 	lotusapi "github.com/filecoin-project/lotus/api"
// 	"github.com/filecoin-project/lotus/chain/types"
// 	"github.com/filecoin-project/specs-actors/actors/builtin/market"
// 	"github.com/ipfs/go-cid"
// )


// type ActiveDeal struct {
// 	Provider address.Address
// }


// func GetChainHead() string {
// 	headers := http.Header{}
// 	addr := "https://node.glif.io/space07/lotus/rpc/v0"
// 	var api lotusapi.FullNodeStruct
// 	closer, err := jsonrpc.NewMergeClient(context.Background(), addr, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
// 	if err != nil {
// 		log.Fatalf("connecting with lotus failed: %s", err)
// 	}
// 	defer closer()


// 	tipset, err := api.ChainHead(context.Background())
// 	if err != nil {
// 		log.Fatalf("calling chain head: %s", err)
// 	}
// 	return tipset
// }

// func GetChainHead2() {
// 	// authToken := "<value found in ~/.lotus/token>"
// 	headers := http.Header{}
// 	addr := "https://node.glif.io/space07/lotus/rpc/v0"

// 	var api lotusapi.FullNodeStruct
// 	closer, err := jsonrpc.NewMergeClient(context.Background(), addr, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
// 	if err != nil {
// 		log.Fatalf("connecting with lotus failed: %s", err)
// 	}
// 	defer closer()

//        // Now you can call any API you're interested in.
// 	tipset, err := api.ChainHead(context.Background())

// 	if err != nil {
// 		log.Fatalf("calling chain head: %s", err)
// 	}


// 	fmt.Printf("==>>\n %+v\n", tipset.Cids()[0])
// 	// fmt.Printf("Current chain head is: %s", tipset.String())\


// 	cidTest, _ := cid.Decode("bafy2bzacecv2qq2ebppqu3mp23oeqwbf2n2xj65ds3il7yxhqhdecuuzty63s")

	

//        // Now you can call any API you're interested in.
// 	messages, err := api.ChainGetBlockMessages(context.Background(), cidTest)




// 	if err != nil {
// 		log.Fatalf("calling chain get message: %s", err)
// 	}

// 	for _, message := range messages.BlsMessages {
// 		// Method 4 is PublishStorageDeals
// 		if (message.Method == 4) {
// 			fmt.Printf("Cid: %+v\n", message.Cid())
			
// 			var params market.PublishStorageDealsParams
// 			err := params.UnmarshalCBOR(bytes.NewReader(message.Params))
// 			if err != nil {
// 					log.Fatalf("UnmarshalCBOR error: %s", err)
// 			}

// 			// Iterate on deals
// 			activeDeals := []ActiveDeal{}
// 			for _, deal := range params.Deals {
// 				proposal := deal.Proposal
// 				// fullJson, _ := json.MarshalIndent(proposal, "", "  ")
// 				// fmt.Printf("proposal: %+v\n", string(fullJson))
// 				if (proposal.VerifiedDeal == true) {
// 					fmt.Printf(" deal.Proposal: %+v\n",  proposal.Provider)
// 					fmt.Printf(" deal.Proposal: %+v\n",  proposal.PieceCID)
// 					activeDeal := ActiveDeal{
// 						Provider: proposal.Provider,
// 					}
// 					activeDeals = append(activeDeals, activeDeal)
// 				}
// 			}


// 			fmt.Printf("activeDeals: %+v\n", activeDeals)

// 		}

// 	}


// 	// test3, _ := api.BeaconGetEntry(context.Background(), abi.ChainEpoch(1070159))
// 	// fmt.Printf("\n ::::>>\n %+v\n", test3)

// 	blockCids, _ := api.ChainGetTipSetByHeight(context.Background(),  abi.ChainEpoch(1070273), types.TipSetKey{})
// 	fmt.Printf("\n blockCids ::::>>\n %+v\n", blockCids)
// 	// for _, cid := range blockCids {
// 	// 	fmt.Printf("tipset: %+v", cid)
// 	// }
	

// 	// fmt.Printf("==>>\n %+v\n", messages)
// 	// fmt.Printf("==>>\n %+v\n", messages.BlsMessages)
// 	// fmt.Printf("==>>\n %+v\n", messages.BlsMessages[0].Method)
// 	// fmt.Printf("Current chain head is: %s", tipset.String())

// }


// func main() {
	
// 	GetChainHead()

	
// }