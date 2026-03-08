package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	billingv1 "github.com/belayhun-arage/billing-service/gen/billing/v1"
)

// paymentExecutor is satisfied by *usecase.ProcessPaymentUsecase and any test mock.
type paymentExecutor interface {
	Execute(ctx context.Context, customerID string, amount int64) error
}

// PaymentHandler implements billingv1.BillingServiceServer.
type PaymentHandler struct {
	billingv1.UnimplementedBillingServiceServer
	usecase paymentExecutor
}

func NewPaymentHandler(u paymentExecutor) *PaymentHandler {
	return &PaymentHandler{usecase: u}
}

func (h *PaymentHandler) ProcessPayment(
	ctx context.Context,
	req *billingv1.ProcessPaymentRequest,
) (*billingv1.ProcessPaymentResponse, error) {

	if req.CustomerId == "" {
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than zero")
	}

	err := h.usecase.Execute(ctx, req.CustomerId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "payment processing failed: %v", err)
	}

	return &billingv1.ProcessPaymentResponse{
		Status: "completed",
	}, nil
}
