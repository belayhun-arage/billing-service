package stripe

import (
	"fmt"

	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// StripeClient wraps the Stripe API for payment processing.
// It implements domain.PaymentProcessor.
type StripeClient struct {
	apiKey string
}

func NewStripeClient(apiKey string) *StripeClient {
	stripe.Key = apiKey
	return &StripeClient{apiKey: apiKey}
}

// Charge creates and immediately confirms a PaymentIntent server-side.
// stripeCustomerID must be a Stripe Customer with a default payment method on file.
// Used for off-session billing (subscriptions, invoices, etc.).
func (s *StripeClient) Charge(stripeCustomerID string, amount int64, currency string) (*domain.ChargeResult, error) {
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amount),
		Currency:           stripe.String(currency),
		Customer:           stripe.String(stripeCustomerID),
		Confirm:            stripe.Bool(true),
		OffSession:         stripe.Bool(true),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe charge failed: %w", err)
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return nil, fmt.Errorf("stripe payment did not succeed: status=%s", pi.Status)
	}

	return &domain.ChargeResult{
		ProviderPaymentID: pi.ID,
		Status:            string(pi.Status),
	}, nil
}

// CreatePaymentIntent creates a PaymentIntent for client-side confirmation via Stripe.js.
// Returns (clientSecret, paymentIntentID, error).
// Use this for checkout flows where the customer enters card details on the frontend.
func (s *StripeClient) CreatePaymentIntent(stripeCustomerID string, amount int64, currency string) (string, string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amount),
		Currency:           stripe.String(currency),
		Customer:           stripe.String(stripeCustomerID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", fmt.Errorf("stripe create payment intent failed: %w", err)
	}

	return pi.ClientSecret, pi.ID, nil
}

// Refund issues a full refund for the given PaymentIntent ID.
func (s *StripeClient) Refund(providerPaymentID string) error {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(providerPaymentID),
	}

	_, err := refund.New(params)
	if err != nil {
		return fmt.Errorf("stripe refund failed: %w", err)
	}

	return nil
}
