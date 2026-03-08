package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type InvoiceHandler struct {
	create      *usecase.CreateInvoiceUsecase
	sendInvoice *usecase.SendInvoiceUsecase
	log         *slog.Logger
}

func NewInvoiceHandler(
	create *usecase.CreateInvoiceUsecase,
	sendInvoice *usecase.SendInvoiceUsecase,
	log *slog.Logger,
) *InvoiceHandler {
	return &InvoiceHandler{create: create, sendInvoice: sendInvoice, log: log}
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

	invoice, err := h.create.Execute(c.Request.Context(), req.CustomerID, req.Amount)
	if err != nil {
		h.log.Error("create invoice failed", "customer_id", req.CustomerID, "amount", req.Amount, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("invoice created", "invoice_id", invoice.ID, "customer_id", invoice.CustomerID)
	c.JSON(http.StatusCreated, invoice)
}

// DownloadPDF streams a PDF of the invoice as a file download.
// GET /invoices/:id/pdf
func (h *InvoiceHandler) DownloadPDF(c *gin.Context) {
	id := c.Param("id")

	pdfData, err := h.sendInvoice.GeneratePDF(c.Request.Context(), id)
	if err != nil {
		h.log.Error("generate PDF failed", "invoice_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="invoice-%.8s.pdf"`, id))
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

// SendByEmail generates a PDF invoice and emails it to the customer.
// POST /invoices/:id/send
func (h *InvoiceHandler) SendByEmail(c *gin.Context) {
	id := c.Param("id")

	h.log.Info("sending invoice by email", "invoice_id", id)

	if err := h.sendInvoice.Send(c.Request.Context(), id); err != nil {
		h.log.Error("send invoice failed", "invoice_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("invoice sent", "invoice_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "invoice sent"})
}
