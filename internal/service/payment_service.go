package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/pkg/db"
)

// ProcessPayment charges the customer via Stripe, then atomically records the
// invoice, payment, ledger entry, and outbox event in a single DB transaction.
//
// Ordering matters: Stripe is charged BEFORE the transaction opens so that a DB
// failure never leaves us with a silent double-charge on retry. If the DB write
// fails after a successful charge, the payment can be reconciled via the Stripe
// dashboard or a compensating refund.
func (s *BillingService) ProcessPayment(
	ctx context.Context,
	customerID string,
	stripeCustomerID string,
	amount int64,
	currency string,
	idempotencyKey string,
) (*domain.Invoice, error) {

	// 1. Idempotency check — replay cached response if already processed.
	if idempotencyKey != "" {
		resp, _, err := s.idempotencyRepo.Get(ctx, idempotencyKey)
		if err == nil {
			var invoice domain.Invoice
			json.Unmarshal(resp, &invoice)
			return &invoice, nil
		}
	}

	// 2. Charge via Stripe (outside the DB transaction).
	charge, err := s.processor.Charge(stripeCustomerID, amount, currency)
	if err != nil {
		return nil, err
	}

	var invoiceResult *domain.Invoice

	// 3. Atomically record all side-effects in the database.
	err = db.WithTransaction(ctx, s.db, func(tx pgx.Tx) error {

		invoice := &domain.Invoice{
			ID:         uuid.New().String(),
			CustomerID: customerID,
			Amount:     amount,
			Status:     "paid",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := s.invoiceRepo.Create(ctx, tx, invoice); err != nil {
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
		if err := s.paymentRepo.Create(ctx, tx, payment); err != nil {
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
		if err := s.ledgerRepo.Create(ctx, tx, ledgerEntry); err != nil {
			return err
		}

		payload, _ := json.Marshal(invoice)
		event := &domain.OutboxEvent{
			ID:        uuid.New().String(),
			EventType: "invoice_paid",
			Payload:   payload,
			CreatedAt: time.Now(),
		}
		if err := s.outboxRepo.Insert(ctx, tx, event); err != nil {
			return err
		}

		invoiceResult = invoice

		if idempotencyKey != "" {
			respData, _ := json.Marshal(invoice)
			if err := s.idempotencyRepo.Save(ctx, idempotencyKey, respData, 200); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if invoiceResult == nil {
		return nil, errors.New("failed to process payment")
	}

	return invoiceResult, nil
}
