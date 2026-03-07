package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type InvoiceHandler struct {
	usecase *usecase.CreateInvoiceUsecase
}

func NewInvoiceHandler(u *usecase.CreateInvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{usecase: u}
}

type CreateInvoiceRequest struct {
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"`
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {

	var req CreateInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	invoice, err := h.usecase.Execute(
		c.Request.Context(),
		req.CustomerID,
		req.Amount,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, invoice)
}