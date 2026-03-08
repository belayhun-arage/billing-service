package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type PaymentHandler struct {
	usecase *usecase.ProcessPaymentUsecase
}

func NewPaymentHandler(u *usecase.ProcessPaymentUsecase) *PaymentHandler {
	return &PaymentHandler{usecase: u}
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

	result, err := h.usecase.Execute(c.Request.Context(), req.CustomerID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment_id": result.PaymentID,
		"invoice_id": result.InvoiceID,
		"status":     "payment processed",
	})
}
