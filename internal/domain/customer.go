package domain

import (
	"context"
	"time"
)

type Customer struct {
	ID               string
	Name             string
	Email            string
	StripeCustomerID string // Stripe Customer ID (set after Stripe registration)
	CreatedAt        time.Time
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id string) (*Customer, error)
}
