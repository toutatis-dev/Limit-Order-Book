# Limit Order Book Engine

A price-time priority matching engine written in Go, exposed via REST API,
with persistent trade logging to PostgreSQL.

## How It Works

Orders are organised into price levels, each containing a FIFO queue.
Buy orders match against the lowest ask, sell orders match against the
highest bid. Unmatched or partially filled orders rest on the book.

## Architecture

- Price levels stored in a map with a sorted slice for best price lookup
- O(1) cancellation via order ID hashmap
- Single-threaded matching engine behind a chi REST API
- Trade logging to PostgreSQL
- CI/CD via GitHub Actions — build, test, vet, then Docker image published to GHCR on successful CI pass

## Endpoints

POST   /order       - Place a limit order

DELETE /order/{id}   - Cancel an order

GET    /book         - Current book state (bids and asks)

## Run

cp .env.example .env   # fill in your DB credentials

docker compose up -d

A pre-built application image is also available at ghcr.io/toutatis-dev/limit-order-book:main


## Test

go test ./...

## Planned

- Prometheus metrics

