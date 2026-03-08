package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type APIKeyHandler struct {
	usecase *usecase.CreateAPIKeyUsecase
	log     *slog.Logger
}

func NewAPIKeyHandler(u *usecase.CreateAPIKeyUsecase, log *slog.Logger) *APIKeyHandler {
	return &APIKeyHandler{usecase: u, log: log}
}

type createAPIKeyRequest struct {
	CustomerID string `json:"customer_id"`
}

// Create issues a new API key + secret for a customer.
// POST /api-keys
func (h *APIKeyHandler) Create(c *gin.Context) {
	var req createAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.Execute(c.Request.Context(), req.CustomerID)
	if err != nil {
		h.log.Error("create api key failed", "customer_id", req.CustomerID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("api key created", "key_id", result.ID, "customer_id", result.CustomerID)
	c.JSON(http.StatusCreated, gin.H{
		"key":         result.Key,
		"secret":      result.Secret,
		"customer_id": result.CustomerID,
		"note":        "Store the secret securely — it will not be shown again.",
	})
}
