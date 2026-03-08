package grpc_test

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	billingv1 "github.com/belayhun-arage/billing-service/gen/billing/v1"
	grpcdelivery "github.com/belayhun-arage/billing-service/internal/delivery/grpc"
	"github.com/belayhun-arage/billing-service/test/mocks"
)

func TestProcessPayment(t *testing.T) {
	tests := []struct {
		name           string
		req            *billingv1.ProcessPaymentRequest
		mockErr        error
		wantCode       codes.Code
		wantStatus     string
		wantCustomerID string
		wantAmount     int64
	}{
		{
			name:           "happy path — valid request processed successfully",
			req:            &billingv1.ProcessPaymentRequest{CustomerId: "cust-123", Amount: 5000},
			wantCode:       codes.OK,
			wantStatus:     "completed",
			wantCustomerID: "cust-123",
			wantAmount:     5000,
		},
		{
			name:     "missing customer_id returns InvalidArgument",
			req:      &billingv1.ProcessPaymentRequest{CustomerId: "", Amount: 5000},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "zero amount returns InvalidArgument",
			req:      &billingv1.ProcessPaymentRequest{CustomerId: "cust-123", Amount: 0},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "negative amount returns InvalidArgument",
			req:      &billingv1.ProcessPaymentRequest{CustomerId: "cust-123", Amount: -100},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "usecase failure returns Internal error",
			req:      &billingv1.ProcessPaymentRequest{CustomerId: "cust-123", Amount: 5000},
			mockErr:  errors.New("db connection lost"),
			wantCode: codes.Internal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mocks.MockPaymentExecutor{Err: tc.mockErr}
			handler := grpcdelivery.NewPaymentHandler(mock)

			resp, err := handler.ProcessPayment(context.Background(), tc.req)

			if tc.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if resp.Status != tc.wantStatus {
					t.Errorf("expected response status %q, got %q", tc.wantStatus, resp.Status)
				}
				if mock.CalledWith.CustomerID != tc.wantCustomerID {
					t.Errorf("usecase called with customer_id %q, want %q", mock.CalledWith.CustomerID, tc.wantCustomerID)
				}
				if mock.CalledWith.Amount != tc.wantAmount {
					t.Errorf("usecase called with amount %d, want %d", mock.CalledWith.Amount, tc.wantAmount)
				}
				return
			}

			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("error is not a gRPC status error: %v", err)
			}
			if st.Code() != tc.wantCode {
				t.Errorf("expected gRPC code %v, got %v", tc.wantCode, st.Code())
			}
		})
	}
}
