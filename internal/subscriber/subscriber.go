package subscriber

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SubscribeToBlocks(webSocketURL string) {
	client, err := ethclient.Dial(webSocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

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
		}
	}


}