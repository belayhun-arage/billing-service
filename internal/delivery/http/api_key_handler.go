package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/internal/usecase"
)

type APIKeyHandler struct {
	create *usecase.CreateAPIKeyUsecase
	revoke *usecase.RevokeAPIKeyUsecase
	log    *slog.Logger
}

func NewAPIKeyHandler(
	create *usecase.CreateAPIKeyUsecase,
	revoke *usecase.RevokeAPIKeyUsecase,
	log *slog.Logger,
) *APIKeyHandler {
	return &APIKeyHandler{create: create, revoke: revoke, log: log}
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

	result, err := h.create.Execute(c.Request.Context(), req.CustomerID)
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

// Revoke immediately invalidates an API key.
// DELETE /api-keys/:key  (protected — requires HMAC auth)
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	key := c.Param("key")
	callerCustomerID := c.GetString("customer_id")

	if err := h.revoke.Execute(c.Request.Context(), key, callerCustomerID); err != nil {
		if errors.Is(err, domain.ErrAPIKeyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("revoke api key failed", "key", key, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("api key revoked", "key", key)
	c.JSON(http.StatusOK, gin.H{"message": "api key revoked"})
}
