package domain_test

import (
	"testing"
	"time"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

const testMerchantID = "merchant-abc-123"

// ── NewCustomer ───────────────────────────────────────────────────────────────

func TestNewCustomer_Valid(t *testing.T) {
	c, err := domain.NewCustomer(testMerchantID, "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ID == "" {
		t.Error("ID must not be empty")
	}
	if c.MerchantID != testMerchantID {
		t.Errorf("MerchantID = %q, want %q", c.MerchantID, testMerchantID)
	}
	if c.Name != "Alice" {
		t.Errorf("Name = %q, want %q", c.Name, "Alice")
	}
	if c.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", c.Email, "alice@example.com")
	}
	if c.CreatedAt.IsZero() {
		t.Error("CreatedAt must not be zero")
	}
}

func TestNewCustomer_NormalizesInput(t *testing.T) {
	c, err := domain.NewCustomer(testMerchantID, "  Alice  ", "  ALICE@Example.COM  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name != "Alice" {
		t.Errorf("Name = %q, want trimmed %q", c.Name, "Alice")
	}
	if c.Email != "alice@example.com" {
		t.Errorf("Email = %q, want lowercased %q", c.Email, "alice@example.com")
	}
}

func TestNewCustomer_InvalidInputs(t *testing.T) {
	tests := []struct {
		desc       string
		merchantID string
		name       string
		email      string
	}{
		{"empty merchant_id", "", "Alice", "alice@example.com"},
		{"empty name", testMerchantID, "", "alice@example.com"},
		{"whitespace name", testMerchantID, "   ", "alice@example.com"},
		{"empty email", testMerchantID, "Alice", ""},
		{"no @ in email", testMerchantID, "Alice", "notanemail"},
		{"no dot in domain", testMerchantID, "Alice", "alice@nodot"},
		{"empty local part", testMerchantID, "Alice", "@example.com"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := domain.NewCustomer(tc.merchantID, tc.name, tc.email)
			if err == nil {
				t.Errorf("NewCustomer(%q, %q, %q): expected error, got nil", tc.merchantID, tc.name, tc.email)
			}
		})
	}
}

// ── NewInvoice ────────────────────────────────────────────────────────────────

func TestNewInvoice_Valid(t *testing.T) {
	inv, err := domain.NewInvoice("cust-123", 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.ID == "" {
		t.Error("ID must not be empty")
	}
	if inv.Currency != "usd" {
		t.Errorf("Currency = %q, want %q", inv.Currency, "usd")
	}
	if inv.Status != "pending" {
		t.Errorf("Status = %q, want %q", inv.Status, "pending")
	}
	if inv.CustomerID != "cust-123" {
		t.Errorf("CustomerID = %q, want %q", inv.CustomerID, "cust-123")
	}
}

func TestNewInvoice_InvalidInputs(t *testing.T) {
	tests := []struct {
		desc       string
		customerID string
		amount     int64
	}{
		{"empty customer_id", "", 5000},
		{"zero amount", "cust-123", 0},
		{"negative amount", "cust-123", -1},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := domain.NewInvoice(tc.customerID, tc.amount)
			if err == nil {
				t.Errorf("NewInvoice(%q, %d): expected error, got nil", tc.customerID, tc.amount)
			}
		})
	}
}

// ── NewPayment ────────────────────────────────────────────────────────────────

func TestNewPayment_Valid(t *testing.T) {
	p, err := domain.NewPayment("inv-123", "cust-123", "pi_abc", 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID == "" {
		t.Error("ID must not be empty")
	}
	if p.Status != "completed" {
		t.Errorf("Status = %q, want %q", p.Status, "completed")
	}
	if p.ProviderPaymentID != "pi_abc" {
		t.Errorf("ProviderPaymentID = %q, want %q", p.ProviderPaymentID, "pi_abc")
	}
}

func TestNewPayment_InvalidInputs(t *testing.T) {
	tests := []struct {
		desc       string
		invoiceID  string
		customerID string
		amount     int64
	}{
		{"empty invoice_id", "", "cust-123", 5000},
		{"empty customer_id", "inv-123", "", 5000},
		{"zero amount", "inv-123", "cust-123", 0},
		{"negative amount", "inv-123", "cust-123", -100},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := domain.NewPayment(tc.invoiceID, tc.customerID, "pi_abc", tc.amount)
			if err == nil {
				t.Errorf("NewPayment(%q, %q, %d): expected error, got nil",
					tc.invoiceID, tc.customerID, tc.amount)
			}
		})
	}
}

// ── NewAPIKey ─────────────────────────────────────────────────────────────────

func TestNewAPIKey_Valid(t *testing.T) {
	k, err := domain.NewAPIKey(testMerchantID, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k.ID == "" {
		t.Error("ID must not be empty")
	}
	if k.MerchantID != testMerchantID {
		t.Errorf("MerchantID = %q, want %q", k.MerchantID, testMerchantID)
	}
	if len(k.Key) < 3 || k.Key[:3] != "bk_" {
		t.Errorf("Key %q must start with 'bk_'", k.Key)
	}
	if len(k.Secret) != 64 {
		t.Errorf("Secret length = %d, want 64 hex chars", len(k.Secret))
	}
	if k.Label != "production" {
		t.Errorf("Label = %q, want %q", k.Label, "production")
	}
	if k.RevokedAt != nil {
		t.Error("RevokedAt must be nil on a new key")
	}
}

func TestNewAPIKey_EmptyLabelAllowed(t *testing.T) {
	k, err := domain.NewAPIKey(testMerchantID, "")
	if err != nil {
		t.Fatalf("expected no error for empty label, got %v", err)
	}
	if k.Label != "" {
		t.Errorf("Label = %q, want empty", k.Label)
	}
}

func TestNewAPIKey_MissingMerchantID(t *testing.T) {
	_, err := domain.NewAPIKey("", "production")
	if err == nil {
		t.Error("expected error when merchant_id is empty, got nil")
	}
}

func TestNewAPIKey_GeneratesUniqueKeys(t *testing.T) {
	k1, _ := domain.NewAPIKey(testMerchantID, "production")
	k2, _ := domain.NewAPIKey(testMerchantID, "production")

	if k1.Key == k2.Key {
		t.Error("two generated keys must be unique")
	}
	if k1.Secret == k2.Secret {
		t.Error("two generated secrets must be unique")
	}
}

// ── APIKey.IsActive ───────────────────────────────────────────────────────────

func TestAPIKey_IsActive(t *testing.T) {
	k, _ := domain.NewAPIKey(testMerchantID, "test")

	if !k.IsActive() {
		t.Error("new key should be active")
	}

	now := time.Now()
	k.RevokedAt = &now

	if k.IsActive() {
		t.Error("key with RevokedAt set should not be active")
	}
}
