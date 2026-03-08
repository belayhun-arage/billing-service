package usecase

import (
	"context"
	"errors"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type CreateCustomerUsecase struct {
	repo CustomerRepository
}

func NewCreateCustomerUsecase(r CustomerRepository) *CreateCustomerUsecase {
	return &CreateCustomerUsecase{repo: r}
}

func (u *CreateCustomerUsecase) Execute(ctx context.Context, name, email string) (*domain.Customer, error) {
	customer, err := domain.NewCustomer(name, email)
	if err != nil {
		return nil, err
	}

	exists, err := u.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("a customer with this email already exists")
	}

	if err := u.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
