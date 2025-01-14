package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum-optimism/optimism/op-celestia/celestia"
	openrpc "github.com/rollkit/celestia-openrpc"
	"github.com/rollkit/celestia-openrpc/types/share"
)

func main() {
	if len(os.Args) < 4 {
		panic("usage: op-celestia <namespace> <eth calldata> <auth token>")
	}

	data, _ := hex.DecodeString(os.Args[2])

	frameRef := celestia.FrameRef{}
	frameRef.UnmarshalBinary(data)

	fmt.Printf("celestia block height: %v; tx index: %v\n", frameRef.BlockHeight, frameRef.TxCommitment)
	fmt.Println("-----------------------------------------")
	client, err := openrpc.NewClient(context.Background(), "http://localhost:26658", os.Args[3])
	if err != nil {
		panic(err)
	}
	nsBytes, err := hex.DecodeString(os.Args[1])
	if err != nil {
		panic(err)
	}
	namespace, err := share.NewBlobNamespaceV0(nsBytes)
	if err != nil {
		panic(err)
	}

	namespacedData, err := client.Blob.GetAll(context.Background(), uint64(frameRef.BlockHeight), []share.Namespace{namespace})
	if err != nil {
		panic(err)
	}
	fmt.Printf("optimism block data on celestia: %x\n", namespacedData[0].Data)
}
