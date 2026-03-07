CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,
    account_id TEXT NOT NULL, -- e.g., customer balance, revenue
    type TEXT NOT NULL, -- debit / credit
    reference_id UUID, -- invoice or payment ID
    amount BIGINT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now()
);