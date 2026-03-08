# Billing Service

A production-grade billing backend API built in Go using Clean Architecture principles.

## Features

- **REST API** — Gin HTTP framework
- **PostgreSQL** — pgx/v5 with connection pooling
- **Transactional operations** — invoices, payments, and ledger entries in a single atomic transaction
- **Idempotent payments** — `Idempotency-Key` header deduplicates requests at the middleware level
- **Outbox Pattern** — reliable event publishing to Kafka via a background worker with `FOR UPDATE SKIP LOCKED`
- **Distributed locks** — PostgreSQL advisory locks for singleton worker processes
- **Clean Architecture** — domain → usecase → delivery, with no framework leakage into business logic

---

## Project Structure

```
billing-service/
├── cmd/api/             # Entry point (main.go)
├── internal/
│   ├── domain/          # Entities and repository interfaces
│   ├── usecase/         # Business logic (no framework dependencies)
│   ├── delivery/http/   # Gin HTTP handlers
│   ├── repository/      # PostgreSQL implementations
│   │   └── postgres/
│   ├── service/         # Orchestration service (full payment flow)
│   ├── worker/          # Background workers (OutboxWorker)
│   └── messaging/       # Event publishers (Kafka)
├── pkg/db/              # Connection pool, transactions, advisory locks
│   └── middleware/      # Idempotency middleware
├── external/stripe/     # Stripe client stub
├── migrations/          # SQL migration files (001–008)
├── legacy/              # Original console billing app (v1, kept for history)
└── .env                 # Environment variables (not committed)
```

---

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL 14+

### Environment

Copy `.env` and set your database URL:

```bash
DATABASE_URL=postgres://postgres:password@localhost:5432/billing_db
```

### Run Migrations

Run the SQL files in `migrations/` in order against your database:

```bash
psql $DATABASE_URL -f migrations/001_create_customers.sql
psql $DATABASE_URL -f migrations/002_create_subscriptions.sql
psql $DATABASE_URL -f migrations/003_create_invoices.sql
psql $DATABASE_URL -f migrations/004_create_payments.sql
psql $DATABASE_URL -f migrations/005_create_ledger_entries.sql
psql $DATABASE_URL -f migrations/006_create_idempotency_keys.sql
psql $DATABASE_URL -f migrations/007_create_outbox_events.sql
psql $DATABASE_URL -f migrations/008_create_indexes.sql
```

### Run the Server

```bash
go run ./cmd/api
```

Server starts on `:8080`.

---

## API Endpoints

### Create Customer

```
POST /customers
Content-Type: application/json

{
  "name": "Belayhun Arage",
  "email": "belayhun@example.com"
}
```

### Create Invoice

```
POST /invoices
Content-Type: application/json

{
  "customer_id": "<uuid>",
  "amount": 5000
}
```

### Process Payment (idempotent)

```
POST /payments
Content-Type: application/json
Idempotency-Key: <unique-key>

{
  "customer_id": "<uuid>",
  "amount": 5000
}
```

Repeating the same request with the same `Idempotency-Key` returns the cached response without re-processing.

---

## Key Design Decisions

### Idempotency Middleware

The `Idempotency-Key` header is checked at the Gin middleware layer. If a matching key exists in the `idempotency_keys` table, the stored response is replayed immediately. Otherwise, the response is captured and persisted after the handler completes.

### Outbox Pattern

Payments write an `outbox_events` row atomically within the same database transaction. A background `OutboxWorker` polls `FOR UPDATE SKIP LOCKED`, publishes to Kafka, and marks events processed — ensuring at-least-once delivery without dual-write risk.

### Advisory Locks

`pkg/db/locks.go` provides `AcquireAdvisoryLock` and `TryAdvisoryLock` using PostgreSQL session-level advisory locks. This prevents multiple instances of the `OutboxWorker` from processing the same events in a multi-replica deployment.

---

## Migrations

| File | Creates |
|------|---------|
| `001` | `customers` |
| `002` | `subscriptions` |
| `003` | `invoices` |
| `004` | `payments` |
| `005` | `ledger_entries` |
| `006` | `idempotency_keys` |
| `007` | `outbox_events` |
| `008` | Indexes |
