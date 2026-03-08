package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type PaymentHandler struct {
	usecase *usecase.ProcessPaymentUsecase
	log     *slog.Logger
}

func NewPaymentHandler(u *usecase.ProcessPaymentUsecase, log *slog.Logger) *PaymentHandler {
	return &PaymentHandler{usecase: u, log: log}
}

type ProcessPaymentRequest struct {
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"`
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {

	var req ProcessPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("processing payment", "customer_id", req.CustomerID, "amount", req.Amount)

	result, err := h.usecase.Execute(c.Request.Context(), req.CustomerID, req.Amount)
	if err != nil {
		h.log.Error("process payment failed", "customer_id", req.CustomerID, "amount", req.Amount, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("payment processed", "payment_id", result.PaymentID, "invoice_id", result.InvoiceID)
	c.JSON(http.StatusCreated, gin.H{
		"payment_id": result.PaymentID,
		"invoice_id": result.InvoiceID,
		"status":     "payment processed",
	})
}
