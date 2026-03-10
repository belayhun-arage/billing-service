CREATE TABLE merchants (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
