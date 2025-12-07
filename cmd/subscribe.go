package cmd

import (
	"context"
	// "fmt"
	"log"

	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/seeques/wallet-tracker/internal/fetcher"
	"github.com/seeques/wallet-tracker/internal/subscriber"
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to new blocks",
	Long:  `Continously listen for and display new blocks as they are added to the blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		addresses := map[common.Address]bool{
			common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"): true,
		}

		client, headers, sub, err := subscriber.Subscribe(webSocketURL)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		// Get chain ID to fetch tx's from address
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case header := <-headers:
				fetcher.TrackWallets(client, header, addresses, chainID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
