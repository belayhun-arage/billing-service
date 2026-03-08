package mocks

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/usecase"
)

// MockPaymentExecutor is a test double for the paymentExecutor interface.
// Set Err to simulate a usecase failure; Result to control the returned IDs.
type MockPaymentExecutor struct {
	Err        error
	Result     *usecase.PaymentResult
	CalledWith struct {
		CustomerID string
		Amount     int64
	}
}

func (m *MockPaymentExecutor) Execute(_ context.Context, customerID string, amount int64) (*usecase.PaymentResult, error) {
	m.CalledWith.CustomerID = customerID
	m.CalledWith.Amount = amount
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Result != nil {
		return m.Result, nil
	}
	return &usecase.PaymentResult{PaymentID: "pay-test", InvoiceID: "inv-test"}, nil
}
