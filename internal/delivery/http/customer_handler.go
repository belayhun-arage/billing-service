package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type CustomerHandler struct {
	usecase *usecase.CreateCustomerUsecase
}

func NewCustomerHandler(u *usecase.CreateCustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: u}
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

	customer, err := h.usecase.Execute(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}
