package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type CustomerHandler struct {
	usecase *usecase.CreateCustomerUsecase
	log     *slog.Logger
}

func NewCustomerHandler(u *usecase.CreateCustomerUsecase, log *slog.Logger) *CustomerHandler {
	return &CustomerHandler{usecase: u, log: log}
}

type CreateCustomerRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("creating customer", "email", req.Email)

	customer, err := h.usecase.Execute(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		h.log.Error("create customer failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("customer created", "customer_id", customer.ID)
	c.JSON(http.StatusCreated, customer)
}
