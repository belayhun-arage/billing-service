package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type MerchantHandler struct {
	usecase *usecase.CreateMerchantUsecase
	log     *slog.Logger
}

func NewMerchantHandler(u *usecase.CreateMerchantUsecase, log *slog.Logger) *MerchantHandler {
	return &MerchantHandler{usecase: u, log: log}
}

type createMerchantRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// Create registers a new merchant (tenant).
// POST /merchants  (public — bootstrapping endpoint)
func (h *MerchantHandler) Create(c *gin.Context) {
	var req createMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, err := h.usecase.Execute(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		h.log.Error("create merchant failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("merchant created", "merchant_id", merchant.ID)
	c.JSON(http.StatusCreated, merchant)
}
