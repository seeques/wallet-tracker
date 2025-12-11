package fetcher

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/seeques/wallet-tracker/storage"
	"github.com/rs/zerolog/log"
)

func TrackWallets(client *ethclient.Client, header *types.Header, addresses map[common.Address]bool, chainID *big.Int, store storage.Storage) error {
	fmt.Printf("New block header: %v\n", header.Hash().Hex())

	block, err := client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		return err
	}

	for _, tx := range block.Transactions() {
		from, err := types.Sender(types.NewPragueSigner(chainID), tx)
		if err != nil {
			return err
		}

		if addresses[from] {
			toAddr := ""
			if tx.To() != nil {
				toAddr = tx.To().Hex()
			}

			log.Info().
    			Str("tx", tx.Hash().Hex()).
    			Str("from", from.Hex()).
    			Str("to", toAddr).
    			Msg("matched transaction")

			tracked := &storage.TrackedTransaction{
				Hash:        tx.Hash().Hex(),
				BlockNumber: header.Number.Uint64(),
				From:        from.Hex(),
				To:          toAddr,
				Value:       tx.Value(),
				Timestamp:   header.Time,
			}
			err := store.SaveTransaction(context.Background(), tracked)

			if err != nil {
				return err
			}
		}

		if tx.To() != nil && addresses[*tx.To()] {

			log.Info().
    			Str("tx", tx.Hash().Hex()).
    			Str("from", from.Hex()).
    			Str("to", tx.To().Hex()).
    			Msg("matched transaction")

			tracked := &storage.TrackedTransaction{
				Hash:        tx.Hash().Hex(),
				BlockNumber: header.Number.Uint64(),
				From:        from.Hex(),
				To:          tx.To().Hex(),
				Value:       tx.Value(),
				Timestamp:   header.Time,
			}
			err := store.SaveTransaction(context.Background(), tracked)

			if err != nil {
				return err
			}
		}
	}
	return nil
}
