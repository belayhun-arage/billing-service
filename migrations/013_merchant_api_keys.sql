-- Bind API keys to a merchant (tenant).
-- Existing keyless rows must be removed or backfilled before applying this migration.
ALTER TABLE api_keys
    ADD COLUMN merchant_id UUID NOT NULL REFERENCES merchants(id);

CREATE INDEX idx_api_keys_merchant_id ON api_keys(merchant_id);
