package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

// type FCOptions struct {
// 	Jsonrpc string
// 	Method string
// 	Id uint64
// 	Params []string
// }

// type ChainHeadResponse struct {
// 	Page   int      `json:"page"`
// 	Fruits []string `json:"fruits"`
// }

// func GetChainHead() {
// 	values := FCOptions {
// 		Jsonrpc: "2.0",
// 		Method: "Filecoin.ChainHead",
// 		Id: 1,
// 		Params: []string{},
// 	}

// 	json_data, err := json.Marshal(values)

// 	resp, err := http.Post(
// 		"https://node.glif.io/space07/lotus/rpc/v0",
// 		"application/json",
// 		bytes.NewBuffer(json_data),
// 	)

//     if err != nil {
//         log.Fatal(err)
//     }

//     defer resp.Body.Close()

//     body, err := ioutil.ReadAll(resp.Body)

//     if err != nil {
//         log.Fatal(err)
//     }

//     fmt.Println(string(body))
// }

type Mytest struct {
	Deals []Mytest2 `json:"Deals"`
	LastName  string `json:"lastName"`
}


type Mytest2 struct {
	Proposal Mytest3 `json:"Proposal"`
}

type Mytest3 struct {
	Label string `json:"Label"`
}


func GetChainHead() {
	// authToken := "<value found in ~/.lotus/token>"
	headers := http.Header{}
	addr := "https://node.glif.io/space07/lotus/rpc/v0"

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), addr, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

       // Now you can call any API you're interested in.
	tipset, err := api.ChainHead(context.Background())

	if err != nil {
		log.Fatalf("calling chain head: %s", err)
	}


	fmt.Printf("==>>\n %+v\n", tipset.Cids()[0])
	// fmt.Printf("Current chain head is: %s", tipset.String())\


	cidTest, _ := cid.Decode("bafy2bzacecv2qq2ebppqu3mp23oeqwbf2n2xj65ds3il7yxhqhdecuuzty63s")
	

       // Now you can call any API you're interested in.
	messages, err := api.ChainGetBlockMessages(context.Background(), cidTest)
	if err != nil {
		log.Fatalf("calling chain get message: %s", err)
	}

	for _, message := range messages.BlsMessages {
		fmt.Printf("message: %+v\n", message.Method)
		if (message.Method == 4) {
			fmt.Printf("YEAAAAAAAAAAAAAAH: %+v\n", message.Method)
			fmt.Printf("YEAAAAAAAAAAAAAAH: %+v\n", message.Cid())
			fmt.Printf("YEAAAAAAAAAAAAAAH: %v\n", message.Params)
			fmt.Printf("YEAAAAAAAAAAAAAAH: %v\n", string(message.Params))

			// dst := []byte{}
			// test, _ := hex.Decode(dst, message.Params)
			// fmt.Printf("YEAAAAAAAAAAAAAAH: %v\n", test)
			// fmt.Printf("YEAAAAAAAAAAAAAAH: %v\n", dst)
			// var p Mytest
			// err := json.Unmarshal(message.Params, &p)
			// if err != nil {
			// 	panic(err)
			// }

			// fmt.Printf("%+v", p)

		}
	}


	// test3, _ := api.BeaconGetEntry(context.Background(), abi.ChainEpoch(1070159))
	// fmt.Printf("\n ::::>>\n %+v\n", test3)

	blockCids, _ := api.ChainGetTipSetByHeight(context.Background(),  abi.ChainEpoch(1070273), types.TipSetKey{})
	for _, cid := range blockCids {
		fmt.Printf("tipset: %+v", cid)
	}
	

	// fmt.Printf("==>>\n %+v\n", messages)
	// fmt.Printf("==>>\n %+v\n", messages.BlsMessages)
	// fmt.Printf("==>>\n %+v\n", messages.BlsMessages[0].Method)
	// fmt.Printf("Current chain head is: %s", tipset.String())

}


func main() {
	
	GetChainHead()

	
}