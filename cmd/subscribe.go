package cmd

import (
	// "context"
	// "fmt"
	// "log"

	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/seeques/wallet-tracker/internal/subscriber"
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to new blocks",
	Long:  `Continously listen for and display new blocks as they are added to the blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		subscriber.SubscribeToBlocks(webSocketURL)
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
