package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByCustomerID(ctx context.Context, customerID string) (*domain.Subscription, error)
}

type CreateSubscriptionUsecase struct {
	repo SubscriptionRepository
}

func NewCreateSubscriptionUsecase(r SubscriptionRepository) *CreateSubscriptionUsecase {
	return &CreateSubscriptionUsecase{repo: r}
}

func (u *CreateSubscriptionUsecase) Execute(ctx context.Context, customerID, plan string) (*domain.Subscription, error) {
	if customerID == "" || plan == "" {
		return nil, errors.New("customerID and plan are required")
	}

	sub := &domain.Subscription{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Plan:       plan,
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	if err := u.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}
