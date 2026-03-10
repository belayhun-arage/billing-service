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
	MerchantID string `json:"merchant_id" binding:"required"`
	Label      string `json:"label"` // optional human-readable identifier, e.g. "production"
}

// Create issues a new merchant-scoped API key + secret.
// POST /api-keys  (public — used to bootstrap the first key for a merchant)
func (h *APIKeyHandler) Create(c *gin.Context) {
	var req createAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.create.Execute(c.Request.Context(), req.MerchantID, req.Label)
	if err != nil {
		h.log.Error("create api key failed", "merchant_id", req.MerchantID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("api key created", "key_id", result.ID, "merchant_id", result.MerchantID, "label", result.Label)
	c.JSON(http.StatusCreated, gin.H{
		"key":         result.Key,
		"secret":      result.Secret,
		"label":       result.Label,
		"merchant_id": result.MerchantID,
		"note":        "Store the secret securely — it will not be shown again.",
	})
}

// Revoke immediately invalidates an API key.
// DELETE /api-keys/:key  (protected — requires a valid HMAC-authenticated request)
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	key := c.Param("key")

	if err := h.revoke.Execute(c.Request.Context(), key); err != nil {
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
