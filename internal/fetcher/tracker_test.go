package fetcher

import (
	"testing"
	"context"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/seeques/wallet-tracker/storage"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setUpSubscription(t *testing.T) (*ethclient.Client, *types.Header, *big.Int, map[common.Address]bool) {
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

	addresses := map[common.Address]bool{common.HexToAddress("0x21a31Ee1afC51d94C2eFcCAa2092aD1028285549"): true}

	return client, header, chainID, addresses
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        t.Skip("DATABASE_URL not set, skipping test")
    }
    
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    
    return pool
}

func TestTrackedWallets(t *testing.T) {
	client, header, chainID, addresses := setUpSubscription(t)

	pool := setupTestDB(t)
	defer pool.Close()

	storage := storage.NewPgStorage(pool)

	err := TrackWallets(client, header, addresses, chainID, storage)
	if err != nil {
		t.Fatal(err)
	}
}