CREATE TABLE api_keys (
    id          UUID PRIMARY KEY,
    key         TEXT UNIQUE NOT NULL,
    secret      TEXT NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers(id),
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    revoked_at  TIMESTAMP
);

CREATE INDEX idx_api_keys_key ON api_keys(key);
