package storage

import (
	"context"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connection string: postgres://tracker:tracker@localhost:5432/wallet_tracker?sslmode=disable

type TrackedTransaction struct {
	Hash        string
	BlockNumber uint64
	From        string
	To          string
	Value       *big.Int
	Timestamp   uint64
}

type Storage interface {
	SaveTransaction(ctx context.Context, tx *TrackedTransaction) error
}

func CreatePool() (*Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), "postgres://tracker:tracker@localhost:5432/wallet_tracker?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
