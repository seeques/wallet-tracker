package fetcher

import (
	"context"
	"fmt"
	"log"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TrackWallets(client *ethclient.Client, header *types.Header, addresses map[common.Address]bool) {
	// Get chain ID to fetch tx's from address
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	fmt.Printf("New block header: %v\n", header.Hash().Hex())

	block, _ := client.BlockByHash(context.Background(), header.Hash())

	for _, tx := range block.Transactions() {
		from, err := types.Sender(types.NewLondonSigner(chainID), tx)
		if err != nil {
			log.Fatalf("Failed to get the sender: %v", err)
		}

		if addresses[from] == true {
			fmt.Printf("Spotted from address %s in tx %s\n", from.Hex(), tx.Hash().Hex())
		}

		if addresses[*tx.To()] == true {
			fmt.Printf("Spotted to address %s in tx %s\n", tx.To().Hex(), tx.Hash().Hex())
		}
	}
}
