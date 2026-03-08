package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/pkg/db"
)

type ProcessPaymentUsecase struct {
	pool        *pgxpool.Pool
	invoiceRepo domain.InvoiceRepository
	paymentRepo domain.PaymentRepository
}

func NewProcessPaymentUsecase(
	pool *pgxpool.Pool,
	invoiceRepo domain.InvoiceRepository,
	paymentRepo domain.PaymentRepository,
) *ProcessPaymentUsecase {
	return &ProcessPaymentUsecase{
		pool:        pool,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
	}
}
// PaymentResult holds the IDs created during payment processing.
type PaymentResult struct {
	PaymentID string
	InvoiceID string
}

func (u *ProcessPaymentUsecase) Execute(
	ctx context.Context,
	customerID string,
	amount int64,
) (*PaymentResult, error) {

	var result PaymentResult

	err := db.WithTransaction(ctx, u.pool, func(tx pgx.Tx) error {

		invoice := &domain.Invoice{
			ID:         uuid.New().String(),
			CustomerID: customerID,
			Amount:     amount,
			Status:     "pending",
			CreatedAt:  time.Now(),
		}

		if err := u.invoiceRepo.Create(ctx, tx, invoice); err != nil {
			return err
		}

		payment := &domain.Payment{
			ID:        uuid.New().String(),
			InvoiceID: invoice.ID,
			Amount:    amount,
			Status:    "completed",
		}

		if err := u.paymentRepo.Create(ctx, tx, payment); err != nil {
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