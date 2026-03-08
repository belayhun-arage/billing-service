package usecase

import (
	"context"

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

	invoice, err := domain.NewInvoice(customerID, amount)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, nil, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}