package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type InvoiceHandler struct {
	usecase *usecase.CreateInvoiceUsecase
	log     *slog.Logger
}

func NewInvoiceHandler(u *usecase.CreateInvoiceUsecase, log *slog.Logger) *InvoiceHandler {
	return &InvoiceHandler{usecase: u, log: log}
}

type CreateInvoiceRequest struct {
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"`
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {

	var req CreateInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("creating invoice", "customer_id", req.CustomerID, "amount", req.Amount)

	invoice, err := h.usecase.Execute(
		c.Request.Context(),
		req.CustomerID,
		req.Amount,
	)

	if err != nil {
		h.log.Error("create invoice failed", "customer_id", req.CustomerID, "amount", req.Amount, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("invoice created", "invoice_id", invoice.ID, "customer_id", invoice.CustomerID)
	c.JSON(http.StatusCreated, invoice)
}