package storage 

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connection string: postgres://tracker:tracker@localhost:5432/wallet_tracker?sslmode=disable

func CreatePool() (*Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), "postgres://tracker:tracker@localhost:5432/wallet_tracker?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
