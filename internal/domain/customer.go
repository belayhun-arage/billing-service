package domain

import "time"

type Customer struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
