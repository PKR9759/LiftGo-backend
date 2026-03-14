# LiftGo-backend

Go REST API backend for LiftGo — a peer-to-peer ride-sharing platform
connecting drivers with empty seats to riders on the same route.

## Tech stack

- Go 1.22
- PostgreSQL
- chi router
- pgx driver
- JWT auth

## Getting started

1. Clone the repo
2. Copy `.env.example` to `.env` and fill in values
3. Run `go run ./cmd/server`

## API docs

See [API.md](./API.md) for all endpoints.