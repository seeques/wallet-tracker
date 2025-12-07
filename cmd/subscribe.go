package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
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
			common.HexToAddress("0x21a31Ee1afC51d94C2eFcCAa2092aD1028285549"): true,
		}

		client, headers, sub, err := subscriber.Subscribe(webSocketURL)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		defer sub.Unsubscribe()

		// Get chain ID to fetch tx's from address
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// create worker channel that will take from receiver routine
		var wg sync.WaitGroup
		worker := make(chan *types.Header, 10)

		// Receiver logic
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(worker) // if we stop, processor need to know that no more workers will arrive
			for {
				select {
				case err := <-sub.Err():
					log.Fatal(err)
					return
				case header := <-headers:
					worker <- header // populate worker
				}
			} 
		}()
		
		// Processor logic
		wg.Add(1)
		go func() {
			defer wg.Done()
			for header := range worker {
				err := fetcher.TrackWallets(client, header, addresses, chainID)
				fmt.Errorf("%v", err)
			}
		}()

		// Wait for goroutines
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
