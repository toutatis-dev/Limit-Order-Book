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
		return nil, failure
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
		return nil, failure
	}
	return &store, nil
}

func (s *Store) WriteTrade(trade engine.Trade) error {
	// write the trade to the postgres store
	_, err := s.Pool.Exec(context.Background(),
		`INSERT INTO trades (taker_id, maker_id, price, quantity, timestamp)
		VALUES ($1, $2, $3, $4, $5);
		`,
		trade.TakerID, trade.MakerID, trade.Price, trade.Quantity, trade.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("Could not write order: %v", err)
	}
	return nil
}

func (s *Store) Close() {
	s.Pool.Close()
}
