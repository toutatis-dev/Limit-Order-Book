package store

import (
	"context"
	"lob/internal/engine"
	"os"
	"testing"
	"time"
)

func TestWriteTrade(t *testing.T) {

	pgUrl := os.Getenv("TEST_DB_URL")
	if pgUrl == "" {
		t.Skip("No Test DB url env var set. Skipping DB integration test.")
	}

	testStore, err := NewStore(pgUrl)
	if err != nil {
		t.Fatalf("Could not initialise the test DB connection: %v", err)
	}

	expectedTime := time.Now().UTC().Truncate(time.Microsecond) //truncate to microsecond to match what postgres will return

	trade := engine.Trade{
		TakerID:   1,
		MakerID:   2,
		Price:     100,
		Quantity:  20,
		Timestamp: expectedTime,
	}

	err = testStore.WriteTrade(trade)
	if err != nil {
		t.Fatalf("Could not write order: %v", err)
	}

	row := testStore.Pool.QueryRow(context.Background(),
		"SELECT taker_id, maker_id, price, quantity, timestamp FROM trades WHERE timestamp=$1",
		expectedTime,
	)

	var tradeReturn engine.Trade

	err = row.Scan(
		&tradeReturn.TakerID,
		&tradeReturn.MakerID,
		&tradeReturn.Price,
		&tradeReturn.Quantity,
		&tradeReturn.Timestamp,
	)
	if err != nil {
		t.Errorf("Failed to scan row: %v", err)
	}

	if trade != tradeReturn {
		t.Errorf("trade and DB do not match. Expected %v, got %v", trade, tradeReturn)
	}

	_, err = testStore.Pool.Exec(context.Background(), "DELETE FROM trades WHERE timestamp=$1", expectedTime)
	if err != nil {
		t.Errorf("Could not delete row from table: %v", err)
	}

}
