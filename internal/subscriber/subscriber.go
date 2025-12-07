package subscriber

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
)

const listenAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

func SubscribeToBlocks(webSocketURL string) {
	client, err := ethclient.Dial(webSocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	// TODO: check logic. defer actually never executes because of the infinite loop below
	defer client.Close()

	// Get chain ID to fetch tx's from address
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	mockAddress := common.HexToAddress(listenAddr)

	headers := make(chan *types.Header)
	// SubscribeNewHead returns subscription to recent headers that are added to our headers channel
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new headers: %v", err)
	}

	for {
		select {
		// error channel from subscription object
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-headers:
			fmt.Printf("New block header: %v\n", header.Hash().Hex())

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatalf("Failed to get block: %v", err)
			}

			for _, tx := range block.Transactions() {
				from, _ := types.Sender(types.NewLondonSigner(chainID), tx)
				// TODO: need to think about quitting because of this
				// if err != nil {
				// 	log.Fatalf("Failed to get the sender: %v", err)
				// }

				if from == mockAddress {
					fmt.Printf("Transaction from %s found in block %d: %s\n", listenAddr, block.Number().Uint64(), tx.Hash().Hex())
				}
				// check nil for contract creation
				if tx.To() != nil && tx.To().Hex() == listenAddr {
					fmt.Printf("Transaction to %s found in block %d: %s\n", listenAddr, block.Number().Uint64(), tx.Hash().Hex())
				}
			}
		}
	}


}