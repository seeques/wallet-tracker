package fetcher

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TrackWallets(client *ethclient.Client, header *types.Header, addresses map[common.Address]bool) error {
	// Get chain ID to fetch tx's from address
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get network ID: %v", err)
	}

	fmt.Printf("New block header: %v\n", header.Hash().Hex())

	block, err := client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		return fmt.Errorf("Failed to fetch block by hash: %v", err)
	}

	for _, tx := range block.Transactions() {
		from, err := types.Sender(types.NewLondonSigner(chainID), tx)
		if err != nil {
			return fmt.Errorf("Failed to get the sender: %v", err)
		}

		if addresses[from] {
			fmt.Printf("Spotted from address %s in tx %s\n", from.Hex(), tx.Hash().Hex())
		}

		if addresses[*tx.To()] && tx.To != nil {
			fmt.Printf("Spotted to address %s in tx %s\n", tx.To().Hex(), tx.Hash().Hex())
		}
	}
	return nil
}
