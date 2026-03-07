-- Quick lookups
CREATE INDEX idx_invoices_customer ON invoices(customer_id);
CREATE INDEX idx_payments_invoice ON payments(invoice_id);
CREATE INDEX idx_subscriptions_next_billing ON subscriptions(next_billing_date);

-- Outbox unprocessed events
CREATE INDEX idx_outbox_unprocessed ON outbox_events(processed);