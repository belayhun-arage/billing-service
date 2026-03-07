package domain

import "time"

type Subscription struct {
	ID         string
	CustomerID string
	Plan       string
	Status     string
	CreatedAt  time.Time
}
