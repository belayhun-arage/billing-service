CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    plan_name TEXT NOT NULL,
    amount BIGINT NOT NULL,
    interval TEXT NOT NULL, -- e.g., monthly, yearly
    next_billing_date TIMESTAMP NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT now()
);