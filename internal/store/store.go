package store

import (
	"context"
	"fmt"

	"lob/internal/engine"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(url string) *Store {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Errorf("Unable to create connection pool: %v", err)
		return &Store{}
	}
	store := Store{
		Pool: pool,
	}
	return &store
}

func (s *Store) WriteTrade(trade engine.Trade) error {
	// write the trade to the postgres store

	return nil
}
