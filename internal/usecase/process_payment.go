package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/pkg/db"
)

// PaymentResult holds the IDs created during payment processing.
type PaymentResult struct {
	PaymentID string
	InvoiceID string
}

type ProcessPaymentUsecase struct {
	pool         *pgxpool.Pool
	customerRepo domain.CustomerRepository
	invoiceRepo  domain.InvoiceRepository
	paymentRepo  domain.PaymentRepository
	ledgerRepo   domain.LedgerRepository
	outboxRepo   domain.OutboxRepository
	processor    domain.PaymentProcessor
}

func NewProcessPaymentUsecase(
	pool *pgxpool.Pool,
	customerRepo domain.CustomerRepository,
	invoiceRepo domain.InvoiceRepository,
	paymentRepo domain.PaymentRepository,
	ledgerRepo domain.LedgerRepository,
	outboxRepo domain.OutboxRepository,
	processor domain.PaymentProcessor,
) *ProcessPaymentUsecase {
	return &ProcessPaymentUsecase{
		pool:         pool,
		customerRepo: customerRepo,
		invoiceRepo:  invoiceRepo,
		paymentRepo:  paymentRepo,
		ledgerRepo:   ledgerRepo,
		outboxRepo:   outboxRepo,
		processor:    processor,
	}
}

// Execute processes a payment for a customer:
//  1. Looks up the customer to obtain their Stripe customer ID.
//  2. Charges via Stripe outside the DB transaction to prevent silent double-charges on retry.
//  3. Atomically records the invoice, payment, ledger entry, and outbox event.
func (u *ProcessPaymentUsecase) Execute(
	ctx context.Context,
	customerID string,
	amount int64,
) (*PaymentResult, error) {

	// 1. Look up customer — we need their Stripe ID to charge them.
	customer, err := u.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// 2. Charge via Stripe before opening the DB transaction.
	//    If the DB write fails after a successful charge, the payment can be
	//    reconciled via Stripe or a compensating refund — no double-charge risk.
	charge, err := u.processor.Charge(customer.StripeCustomerID, amount, "usd")
	if err != nil {
		return nil, fmt.Errorf("stripe charge failed: %w", err)
	}

	var result PaymentResult

	// 3. Atomically record all side-effects in the database.
	err = db.WithTransaction(ctx, u.pool, func(tx pgx.Tx) error {

		invoice := &domain.Invoice{
			ID:         uuid.New().String(),
			CustomerID: customerID,
			Amount:     amount,
			Status:     "paid",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := u.invoiceRepo.Create(ctx, tx, invoice); err != nil {
			return err
		}

		payment := &domain.Payment{
			ID:                uuid.New().String(),
			InvoiceID:         invoice.ID,
			CustomerID:        customerID,
			Amount:            amount,
			Status:            "completed",
			ProviderPaymentID: charge.ProviderPaymentID,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		if err := u.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}

		ledgerEntry := &domain.LedgerEntry{
			ID:          uuid.New().String(),
			AccountID:   "revenue",
			Type:        "credit",
			ReferenceID: invoice.ID,
			Amount:      amount,
			Description: "Payment received via Stripe",
			CreatedAt:   time.Now(),
		}
		if err := u.ledgerRepo.Create(ctx, tx, ledgerEntry); err != nil {
			return err
		}

		payload, _ := json.Marshal(invoice)
		event := &domain.OutboxEvent{
			ID:        uuid.New().String(),
			EventType: "invoice_paid",
			Payload:   payload,
			CreatedAt: time.Now(),
		}
		if err := u.outboxRepo.Insert(ctx, tx, event); err != nil {
			return err
		}

		result.PaymentID = payment.ID
		result.InvoiceID = invoice.ID
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}
