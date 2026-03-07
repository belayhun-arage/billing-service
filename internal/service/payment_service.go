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

func (s *BillingService) ProcessPayment(
	ctx context.Context,
	customerID string,
	amount int64,
	idempotencyKey string,
) (*domain.Invoice, error) {

	// Check Idempotency first
	if idempotencyKey != "" {
		resp, _, err := s.idempotencyRepo.Get(ctx, idempotencyKey)
		if err == nil {
			// Already processed
			var invoice domain.Invoice
			json.Unmarshal(resp, &invoice)
			return &invoice, nil
		}
	}

	var invoiceResult *domain.Invoice

	err := db.WithTransaction(ctx, s.db, func(tx pgx.Tx) error {

		// 1️⃣ Create Invoice
		invoice := &domain.Invoice{
			ID:         uuid.New().String(),
			CustomerID: customerID,
			Amount:     amount,
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := s.invoiceRepo.Create(ctx, tx, invoice); err != nil {
			return err
		}

		// 2️⃣ Process Payment
		payment := &domain.Payment{
			ID:        uuid.New().String(),
			InvoiceID: invoice.ID,
			CustomerID: customerID,
			Amount:    amount,
			Status:    "completed",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}

		// 3️⃣ Ledger Entry
		ledgerEntry := &domain.LedgerEntry{
			ID:        uuid.New().String(),
			AccountID: "revenue", // Example, could be dynamic
			Type:      "credit",
			ReferenceID: invoice.ID,
			Amount:    amount,
			Description: "Payment received",
			CreatedAt: time.Now(),
		}
		if err := s.ledgerRepo.Create(ctx, tx, ledgerEntry); err != nil {
			return err
		}

		// 4️⃣ Outbox Event
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

		invoice.Status = "paid"
		invoiceResult = invoice

		// 5️⃣ Save idempotency key (optional)
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
		return nil, errors.New("failed to process invoice")
	}

	return invoiceResult, nil
}