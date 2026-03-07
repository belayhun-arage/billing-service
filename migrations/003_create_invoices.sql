CREATE TABLE invoices (
    id UUID PRIMARY KEY,
    subscription_id UUID REFERENCES subscriptions(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    amount BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, paid, failed
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);