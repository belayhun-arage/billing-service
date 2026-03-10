-- Scope customers to a merchant (tenant).
-- The global email uniqueness constraint is replaced by per-merchant uniqueness.
ALTER TABLE customers
    DROP CONSTRAINT IF EXISTS customers_email_key;

ALTER TABLE customers
    ADD COLUMN merchant_id UUID NOT NULL REFERENCES merchants(id);

ALTER TABLE customers
    ADD CONSTRAINT customers_merchant_email_key UNIQUE (merchant_id, email);
