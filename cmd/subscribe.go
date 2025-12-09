package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/ethclient"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/seeques/wallet-tracker/internal/fetcher"
	"github.com/seeques/wallet-tracker/internal/subscriber"
	"github.com/seeques/wallet-tracker/internal/config"
	"github.com/seeques/wallet-tracker/storage"
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to new blocks",
	Long:  `Continously listen to and display new blocks as they are added to the blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		config := config.LoadConfig()

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

		// create pool
		pool, err := storage.CreatePool()
		if err != nil {
			fmt.Printf("Unable to connect to database: %v\n", err)
		}
		defer pool.Close()

		storage := storage.NewPgStorage(pool)

		// create worker channel that will take from receiver routine
		var wg sync.WaitGroup
		worker := make(chan *types.Header, 10)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())

		// Receiver logic
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			defer close(worker) // if we stop, processor needs to know that no more workers will arrive
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-sub.Err():
					fmt.Printf("%v", err)
					return
				case header := <-headers:
					worker <- header // populate worker
				}
			}
		}(ctx)

		// Processor logic
		wg.Add(1)
		go func() {
			defer wg.Done()
			for header := range worker {
				err := fetcher.TrackWallets(client, header, config.Addresses, chainID, storage)
				if err != nil {
					fmt.Printf("%v", err)
				}
			}
		}()

		// wait for signal to fire
		// if websocket drops, the wa.Wait() unblocks and program exits
		go func() {
			<-sigs
			cancel()
		}()

		// Wait for goroutines
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
