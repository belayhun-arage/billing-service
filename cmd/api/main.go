package main

import (
	httpdelivery "github.com/belayhun-arage/billing-service/internal/delivery/http"
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

	customerHandler := httpdelivery.NewCustomerHandler(createCustomerUC)
	invoiceHandler := httpdelivery.NewInvoiceHandler(createInvoiceUC)
	paymentHandler := httpdelivery.NewPaymentHandler(processPaymentUC)

	r := gin.Default()

	r.Use(
		middleware.IdempotencyMiddleware(idempotencyRepo),
	)

	r.POST("/customers", customerHandler.CreateCustomer)
	r.POST("/invoices", invoiceHandler.CreateInvoice)
	r.POST("/payments", paymentHandler.ProcessPayment)
	r.Run(":8080")
}
