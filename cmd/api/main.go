package main

import (
	"github.com/belayhun-arage/billing-service/internal/delivery/http"
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

	invoiceRepo := postgres.NewInvoiceRepository(pool)
	paymentRepo := postgres.NewPaymentRepository(pool)
	idempotencyRepo := postgres.NewIdempotencyRepository(pool)

	createInvoiceUC := usecase.NewCreateInvoiceUsecase(invoiceRepo)
	processPaymentUC := usecase.NewProcessPaymentUsecase(pool, invoiceRepo, paymentRepo)

	invoiceHandler := http.NewInvoiceHandler(createInvoiceUC)
	paymentHandler := http.NewPaymentHandler(processPaymentUC)

	r := gin.Default()

	r.Use(
		middleware.IdempotencyMiddleware(idempotencyRepo),
	)

	r.POST("/invoices", invoiceHandler.CreateInvoice)
	r.POST("/payments", paymentHandler.ProcessPayment)
	r.Run(":8080")
}
