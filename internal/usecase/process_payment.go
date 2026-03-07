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
func (u *ProcessPaymentUsecase) Execute(
	ctx context.Context,
	customerID string,
	amount int64,
) error {

	return db.WithTransaction(ctx, u.pool, func(tx pgx.Tx) error {

		invoice := &domain.Invoice{
			ID:         uuid.New().String(),
			CustomerID: customerID,
			Amount:     amount,
			Status:     "pending",
			CreatedAt:  time.Now(),
		}

		err := u.invoiceRepo.Create(ctx, tx, invoice)
		if err != nil {
			return err
		}

		payment := &domain.Payment{
			ID:        uuid.New().String(),
			InvoiceID: invoice.ID,
			Amount:    amount,
			Status:    "completed",
		}

		err = u.paymentRepo.Create(ctx, tx, payment)
		if err != nil {
			return err
		}

		return nil
	})
}