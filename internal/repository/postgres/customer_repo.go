package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type CustomerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
	INSERT INTO customers (id, name, email, created_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, customer.ID, customer.Name, customer.Email, customer.CreatedAt)
	return err
}

func (r *CustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `
	SELECT id, name, email, stripe_customer_id, created_at
	FROM customers
	WHERE id = $1
	`
	var c domain.Customer
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Email, &c.StripeCustomerID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	query := `
	SELECT id, name, email, created_at
	FROM customers
	WHERE email = $1
	`
	var c domain.Customer
	err := r.db.QueryRow(ctx, query, email).Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (r *CustomerRepository) UpdatedAt(ctx context.Context, id string, updatedAt time.Time) error {
	_, err := r.db.Exec(ctx, `UPDATE customers SET updated_at = $1 WHERE id = $2`, updatedAt, id)
	return err
}
