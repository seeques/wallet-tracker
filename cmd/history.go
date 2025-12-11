package cmd

import (
	"context"

	"github.com/spf13/cobra"
	// "github.com/seeques/wallet-tracker/internal/config"
	"github.com/seeques/wallet-tracker/storage"
	"github.com/rs/zerolog/log"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Read transaction values sent from/to tracked wallets",
	Long:  `Read transaction values sent from or to tracked wallet address by specifying the address with a flag`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Info().Msg("Need to provide an address")
			return
		}

		pool, err := storage.CreatePool()
		if err != nil {
			log.Error().Err(err).Msg("Unable to connect to database")
		}
		defer pool.Close()

		store := storage.NewPgStorage(pool)

		txs, err := store.GetByAddress(context.Background(), args[0])
		if err != nil {
			log.Error().Err(err).Msg("Error scanning tracked wallets")
		}

		for i := range txs {
			tx := &txs[i]

			value := tx.Value

			valueU64 := value.Uint64()
			log.Info().Uint64("value", valueU64).Msg("")
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
