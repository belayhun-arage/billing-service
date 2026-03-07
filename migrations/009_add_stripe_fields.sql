-- Add Stripe provider_payment_id to payments for reconciliation
ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS provider_payment_id TEXT;

-- Add Stripe customer ID to customers so we can charge them server-side
ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;
