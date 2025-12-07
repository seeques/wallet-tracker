package subscriber

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Subscribe(webSocketURL string) (*ethclient.Client, chan *types.Header, ethereum.Subscription, error) {
	client, err := ethclient.Dial(webSocketURL)
	if err != nil {
		return nil, nil, nil, err
	}

	headers := make(chan *types.Header)
	// SubscribeNewHead returns subscription to recent headers that are added to our headers channel
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return nil, nil, nil, err
	}

	return client, headers, sub, nil
}
