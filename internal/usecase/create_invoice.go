package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/belayhun-arage/billing-service/internal/domain"
)

type CreateInvoiceUsecase struct {
	repo domain.InvoiceRepository
}

func NewCreateInvoiceUsecase(r domain.InvoiceRepository) *CreateInvoiceUsecase {
	return &CreateInvoiceUsecase{repo: r}
}

func (u *CreateInvoiceUsecase) Execute(
	ctx context.Context,
	customerID string,
	amount int64,
) (*domain.Invoice, error) {

	invoice := &domain.Invoice{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Amount:     amount,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	err := u.repo.Create(ctx, nil, invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}