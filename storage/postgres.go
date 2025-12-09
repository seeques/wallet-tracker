package storage

import (
	"context"
	"math/big"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seeques/wallet-tracker/internal/config"
)

type TrackedTransaction struct {
	Hash        string
	BlockNumber uint64
	From        string
	To          string
	Value       *big.Int
	Timestamp   uint64
}

type PostgresStorage struct {
	pool *pgxpool.Pool
}

type Storage interface {
	SaveTransaction(ctx context.Context, tx *TrackedTransaction) error
}

func (s *PostgresStorage) SaveTransaction(ctx context.Context, tx *TrackedTransaction) error {
	sql := `INSERT INTO transactions (hash, block_number, from_address, to_address, val, tmstp) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.pool.Exec(ctx, sql, tx.Hash, tx.BlockNumber, tx.From, tx.To, tx.Value, tx.Timestamp)
	return err
}

// Create a constructor function
func NewPgStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func CreatePool() (*pgxpool.Pool, error) {
	config := config.LoadConfig()
	conn, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
