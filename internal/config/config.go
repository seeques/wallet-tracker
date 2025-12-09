package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	ETH_RPC_URL string
	Addresses []string
	DatabaseURL string
}

func loadConfig() Config {
	godotenv.Load()

	return Config{
		ETH_RPC_URL: os.Getenv("ETH_RPC_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Addresses: parseAddresses(os.Getenv("WATCHED_ADDRESSES"))
	}
}

func parseAddresses(raw string) map[common.Address]bool {
	result := make(map[commod.Address]bool)
	if raw == "" {
		return result
	}

	for _, addr := range strings.Split(raw, ",") {
		result[common.HexToAddress(addr)] = true
	}

	return result
}