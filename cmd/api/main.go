package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	billingv1 "github.com/belayhun-arage/billing-service/gen/billing/v1"
	httpdelivery "github.com/belayhun-arage/billing-service/internal/delivery/http"
	grpcdelivery "github.com/belayhun-arage/billing-service/internal/delivery/grpc"
	"github.com/belayhun-arage/billing-service/internal/repository/postgres"
	"github.com/belayhun-arage/billing-service/internal/usecase"
	"github.com/belayhun-arage/billing-service/pkg/db"
	"github.com/belayhun-arage/billing-service/pkg/db/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	pool, err := db.NewPostgresPool()
	if err != nil {
		panic(err)
	}

	customerRepo := postgres.NewCustomerRepository(pool)
	invoiceRepo := postgres.NewInvoiceRepository(pool)
	paymentRepo := postgres.NewPaymentRepository(pool)
	idempotencyRepo := postgres.NewIdempotencyRepository(pool)

	createCustomerUC := usecase.NewCreateCustomerUsecase(customerRepo)
	createInvoiceUC := usecase.NewCreateInvoiceUsecase(invoiceRepo)
	processPaymentUC := usecase.NewProcessPaymentUsecase(pool, invoiceRepo, paymentRepo)

	// --- HTTP server ---
	customerHandler := httpdelivery.NewCustomerHandler(createCustomerUC)
	invoiceHandler := httpdelivery.NewInvoiceHandler(createInvoiceUC)
	paymentHandler := httpdelivery.NewPaymentHandler(processPaymentUC)

	r := gin.Default()
	r.Use(middleware.IdempotencyMiddleware(idempotencyRepo))
	r.POST("/customers", customerHandler.CreateCustomer)
	r.POST("/invoices", invoiceHandler.CreateInvoice)
	r.POST("/payments", paymentHandler.ProcessPayment)

	// --- gRPC server ---
	grpcPaymentHandler := grpcdelivery.NewPaymentHandler(processPaymentUC)

	grpcServer := grpc.NewServer()
	billingv1.RegisterBillingServiceServer(grpcServer, grpcPaymentHandler)

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	// Run both servers concurrently. If either fails, log and exit.
	errCh := make(chan error, 2)

	go func() {
		log.Println("gRPC server listening on :9090")
		errCh <- grpcServer.Serve(lis)
	}()

	go func() {
		log.Println("HTTP server listening on :8080")
		errCh <- r.Run(":8080")
	}()

	log.Fatal(<-errCh)
}
