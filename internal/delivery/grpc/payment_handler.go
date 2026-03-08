package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	billingv1 "github.com/belayhun-arage/billing-service/gen/billing/v1"
	"github.com/belayhun-arage/billing-service/internal/usecase"
)

// paymentExecutor is satisfied by *usecase.ProcessPaymentUsecase and any test mock.
type paymentExecutor interface {
	Execute(ctx context.Context, customerID string, amount int64) (*usecase.PaymentResult, error)
}

// PaymentHandler implements billingv1.BillingServiceServer.
type PaymentHandler struct {
	billingv1.UnimplementedBillingServiceServer
	usecase paymentExecutor
	log     *slog.Logger
}

func NewPaymentHandler(u paymentExecutor, log *slog.Logger) *PaymentHandler {
	return &PaymentHandler{usecase: u, log: log}
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

	h.log.Info("processing payment", "customer_id", req.CustomerId, "amount", req.Amount)

	result, err := h.usecase.Execute(ctx, req.CustomerId, req.Amount)
	if err != nil {
		h.log.Error("process payment failed", "customer_id", req.CustomerId, "amount", req.Amount, "error", err)
		return nil, status.Errorf(codes.Internal, "payment processing failed: %v", err)
	}

	h.log.Info("payment processed", "payment_id", result.PaymentID, "invoice_id", result.InvoiceID)
	return &billingv1.ProcessPaymentResponse{
		PaymentId: result.PaymentID,
		InvoiceId: result.InvoiceID,
		Status:    "completed",
	}, nil
}
