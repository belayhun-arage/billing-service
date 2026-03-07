CREATE TABLE idempotency_keys (
    key TEXT PRIMARY KEY,
    response JSONB,
    status_code INT,
    created_at TIMESTAMP DEFAULT now()
);