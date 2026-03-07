package domain

import "time"

type Customer struct {
	ID               string
	Name             string
	Email            string
	StripeCustomerID string // Stripe Customer ID (set after Stripe registration)
	CreatedAt        time.Time
}
