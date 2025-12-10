package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/seeques/wallet-tracker/internal/config"
	"github.com/seeques/wallet-tracker/storage"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Read transaction values sent from/to tracked wallets",
	Long:  `Read transaction values sent from or to tracked wallet address by specifying the address with a flag`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Need to provide an address")
			return
		}

		pool, err := storage.CreatePool()
		if err != nil {
			fmt.Printf("Unable to connect to database, %v\n", err)
		}
		defer pool.Close()

		store := storage.NewPgStorage(pool)

		txs, err := store.GetByAddress(context.Background(), args[0])
		if err != nil {
			fmt.Printf("Error scanning tracked wallet, %v\n", err)
		}

		for i := range txs {
			tx := &txs[i]

			value := tx.Value
			fmt.Printf("Value: %d\n", value)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
