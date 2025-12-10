# Wallet Tracker

A real-time Ethereum wallet tracker that monitors blockchain transactions for specified addresses and stores them in PostgreSQL.

## What It Does

- Subscribes to new Ethereum blocks via WebSocket
- Filters transactions involving watched addresses (as sender or recipient)
- Stores matching transactions in PostgreSQL
- Provides CLI to query transaction history
- Handles graceful shutdown on Ctrl+C

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- An Ethereum RPC endpoint with WebSocket support (Alchemy, Infura, etc.)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/seeques/wallet-tracker
cd wallet-tracker
```

2. Start PostgreSQL:
```bash
docker compose up -d
```

3. Install migration tool:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

4. Run migrations:
```bash
migrate -path migrations -database "postgres://tracker:tracker@localhost:5433/wallet_tracker?sslmode=disable" up
```

5. Create `.env` file:
```env
ETH_RPC_URL=wss://eth-mainnet.g.alchemy.com/v2/your-api-key
DATABASE_URL=postgres://tracker:tracker@localhost:5433/wallet_tracker?sslmode=disable
WATCHED_ADDRESSES=0xabc123...,0xdef456...
```

6. Build and run the tracker:
```bash
go build
./wallet-tracker subscribe
```

## Configuration

All configuration is via environment variables (or `.env` file):

| Variable | Description | Example |
|----------|-------------|---------|
| `ETH_RPC_URL` | WebSocket RPC endpoint | `wss://eth-mainnet.g.alchemy.com/v2/xxx` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://tracker:tracker@localhost:5433/wallet_tracker?sslmode=disable` |
| `WATCHED_ADDRESSES` | Comma-separated addresses to track | `0xabc...,0xdef...` |

## Commands

### Subscribe to blocks
```bash
wallet-tracker subscribe
```
Continuously monitors the blockchain for transactions involving watched addresses. Press Ctrl+C to stop gracefully.

### Query history
```bash
wallet-tracker history 0xabc123...
```
Shows all stored transactions for a given address.

## Project Structure

```
wallet-tracker/
├── cmd/
│   ├── root.go           # Cobra root command
│   ├── subscribe.go      # Block subscription command with receiver and processor logic
│   └── history.go        # History query command
├── internal/
│   ├── config/
│   │   └── config.go     # Configuration loading
│   ├── fetcher/
│   │   └── fetcher.go    # Wallet tracking logic
│   └── subscriber/
│       └── subscriber.go # WebSocket subscription logic
├── storage/
│   └── postgres.go       # Database operations
├── migrations/
│   └── 000001_create_transactions.up.sql
├── docker-compose.yml
├── .env
├── go.mod
└── README.md
```

## Architecture

```
┌─────────────────┐
│  Ethereum RPC   │
│   (WebSocket)   │
└────────┬────────┘
         │ new block headers
         ▼
┌─────────────────┐
│    Receiver     │
│   Goroutine     │
└────────┬────────┘
         │ buffered channel
         ▼
┌─────────────────┐
│   Processor     │
│   Goroutine     │
└────────┬────────┘
         │ matching transactions
         ▼
┌─────────────────┐
│   PostgreSQL    │
│    Storage      │
└─────────────────┘
```

**Receiver goroutine:** Listens for new block headers, passes them to processor via buffered channel.

**Processor goroutine:** Fetches full block, iterates transactions, filters by watched addresses, saves matches to database.

**Graceful shutdown:** Ctrl+C triggers context cancellation, goroutines finish current work, then exit cleanly.

## Database Schema

```sql
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    hash VARCHAR(66) UNIQUE NOT NULL,
    block_number BIGINT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    val NUMERIC NOT NULL,
    tmstp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```