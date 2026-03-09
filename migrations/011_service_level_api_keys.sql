-- Convert api_keys from customer-scoped to service-level credentials.
-- Drops the customer_id FK and adds an optional label for identifying keys.
ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS api_keys_customer_id_fkey;
ALTER TABLE api_keys DROP COLUMN IF EXISTS customer_id;
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS label TEXT NOT NULL DEFAULT '';
