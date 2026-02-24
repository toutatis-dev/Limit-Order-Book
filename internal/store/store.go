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

func NewStore(url string) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		failure := fmt.Errorf("Unable to create connection pool: %v", err)
		return &Store{}, failure
	}
	store := Store{
		Pool: pool,
	}

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS trades (
		id          BIGSERIAL PRIMARY KEY,
		taker_id    BIGINT NOT NULL,
		maker_id    BIGINT NOT NULL,
		price       BIGINT NOT NULL,
		quantity    BIGINT NOT NULL,
		timestamp   TIMESTAMPTZ NOT NULL
		)`,
	)
	if err != nil {
		failure := fmt.Errorf("Could not create table: %v", err)
		return &Store{}, failure
	}
	return &store, nil
}

func (s *Store) WriteTrade(trade engine.Trade) error {
	// write the trade to the postgres store

	return nil
}
