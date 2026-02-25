# Limit Order Book Engine

A price-time priority matching engine written in Go, exposed via REST API,
with persistent trade logging to PostgreSQL and real-time observability via
Prometheus and Grafana.

## How It Works

Orders are organised into price levels, each containing a FIFO queue implemented
as a doubly linked list. Buy orders match against the lowest ask, sell orders
match against the highest bid. Unmatched or partially filled orders rest on the
book until matched or cancelled.

## Architecture

- Price levels stored in a sorted slice with binary search insertion
- O(1) order lookup and cancellation via ID hashmap
- Doubly linked list per price level for O(1) head removal
- Mutex-locked engine for safe concurrent access
- Interface-based dependency injection throughout for testability
- CI/CD via GitHub Actions — build, test, vet on push, Docker image published to GHCR on successful CI

## Dashboard

![Grafana Dashboard](Docs/dashboard.png)

## Tech Stack

- Go 1.25
- chi — HTTP router
- pgx — PostgreSQL driver
- Prometheus + Grafana — metrics and observability
- Docker + docker-compose — containerisation
- GitHub Actions — CI/CD

## Endpoints

POST   /order       - Place a limit order

DELETE /order/{id}   - Cancel an order

GET    /book         - Current book state (bids and asks)

GET    /metrics      - Prometheus metrics

## Run
```bash
cp .env.example .env   # fill in DB credentials
docker compose up -d
```

A pre-built image is available at `ghcr.io/toutatis-dev/limit-order-book:main`

## Test
```bash
go test ./...
# Integration tests require TEST_DB_URL to be set
TEST_DB_URL="postgres://user:pass@localhost:5432/lob" go test ./internal/store/...
```

