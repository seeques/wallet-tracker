package fetcher

import (
	"testing"
	"context"
	"math/big"
	"os"

	"github.com/seeques/wallet-tracker/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/seeques/wallet-tracker/storage"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func setUpSubscription(t *testing.T) (*ethclient.Client, *types.Header, *big.Int, map[common.Address]bool) {
	config := config.LoadConfig()
	rpcURL := os.Getenv("ETH_RPC_URL")

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatal(err)
	}

	block := new(big.Int).SetInt64(23974935)

	header, err := client.HeaderByNumber(context.Background(), block)
	if err != nil {
		t.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return client, header, chainID, config.Addresses
}

func TestTrackedWallets(t *testing.T) {
	client, header, chainID, addresses := setUpSubscription(t)

	pool, err := storage.CreatePool()
	if err != nil {
			t.Fatal(err)
		}
	defer pool.Close()

	storage := storage.NewPgStorage(pool)

	err = TrackWallets(client, header, addresses, chainID, storage)
	if err != nil {
		t.Fatal(err)
	}
}