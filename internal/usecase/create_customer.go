package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

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
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	exists, err := u.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("a customer with this email already exists")
	}

	customer := &domain.Customer{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
