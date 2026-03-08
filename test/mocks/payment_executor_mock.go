package mocks

import "context"

// MockPaymentExecutor is a test double for the paymentExecutor interface.
// Set Err to simulate a usecase failure.
type MockPaymentExecutor struct {
	Err            error
	CalledWith     struct {
		CustomerID string
		Amount     int64
	}
}

func (m *MockPaymentExecutor) Execute(_ context.Context, customerID string, amount int64) error {
	m.CalledWith.CustomerID = customerID
	m.CalledWith.Amount = amount
	return m.Err
}
