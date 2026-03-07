# Billing Service (Go)  

A **robust billing backend API** built in **Go** using **Clean Architecture** principles.  
Supports **PostgreSQL (pgx)**, **transactional operations**, **idempotent payments**, **Outbox pattern for reliable events**, and **distributed worker locks** for high-concurrency fintech systems.  

---

## Features

### Core Features
- REST API built with [Gin](https://github.com/gin-gonic/gin)
- PostgreSQL integration using [pgx](https://github.com/jackc/pgx)
- Domain-driven design & **Clean Architecture**
- Transactional operations for **invoices**, **payments**, and **ledger entries**
- Idempotent payment API using `Idempotency-Key` header
- **Outbox Pattern** for reliable event publishing to message brokers (Kafka, NATS, etc.)
- Distributed locks for billing workers using **PostgreSQL advisory locks**
- Structured logging with [zerolog](https://github.com/rs/zerolog)
- Environment-based configuration using [godotenv](https://github.com/joho/godotenv)

---

## Project Structure
billing-service/
├── cmd/
│ └── api/
│ └── main.go # Entry point
├── internal/
│ ├── domain/ # Core entities and interfaces
│ ├── usecase/ # Application/business logic
│ ├── delivery/http/ # HTTP handlers
│ ├── repository/postgres/ # Database implementation
│ ├── worker/ # Background workers (billing, outbox)
│ └── messaging/ # Event publishers (Kafka, NATS, etc.)
├── pkg/
│ └── db/ # Database connection, transactions, locks
├── migrations/ # SQL migrations
├── .env # Environment variables
└── go.mod