/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/seeques/wallet-tracker/internal/config"
	"github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wallet-tracker",
	Short: "A brief description of your application",
	Long: ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	},
}

var webSocketURL string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config := config.LoadConfig()
	rootCmd.PersistentFlags().StringVar(&webSocketURL, "socket", config.ETH_RPC_URL, "WebSocket URL")
}
