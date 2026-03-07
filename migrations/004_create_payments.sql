CREATE TABLE payments (
    id UUID PRIMARY KEY,
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    amount BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, completed, failed
    method TEXT, -- e.g., card, bank transfer
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);