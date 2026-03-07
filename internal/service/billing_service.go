package service

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/repository/postgres"
)

type BillingService struct {
	db             *pgxpool.Pool
	invoiceRepo    *postgres.InvoiceRepository
	paymentRepo    *postgres.PaymentRepository
	ledgerRepo     *postgres.LedgerRepository
	outboxRepo     *postgres.OutboxRepository
	idempotencyRepo *postgres.IdempotencyRepository
}

func NewBillingService(
	db *pgxpool.Pool,
	invoiceRepo *postgres.InvoiceRepository,
	paymentRepo *postgres.PaymentRepository,
	ledgerRepo *postgres.LedgerRepository,
	outboxRepo *postgres.OutboxRepository,
	idempotencyRepo *postgres.IdempotencyRepository,
) *BillingService {
	return &BillingService{
		db:             db,
		invoiceRepo:    invoiceRepo,
		paymentRepo:    paymentRepo,
		ledgerRepo:     ledgerRepo,
		outboxRepo:     outboxRepo,
		idempotencyRepo: idempotencyRepo,
	}
}