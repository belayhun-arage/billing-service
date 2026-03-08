package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

const timestampTolerance = 5 * time.Minute

// HMACAuth returns a Gin middleware that authenticates requests using an
// API key + HMAC-SHA256 request signature.
//
// Required headers:
//
//	X-API-Key   — the public key identifier  (e.g. bk_a3f1...)
//	X-Timestamp — Unix timestamp in seconds  (e.g. 1710000000)
//	X-Signature — hex(HMAC-SHA256(secret, message))
//
// Signed message:
//
//	METHOD\nPATH\nTIMESTAMP\nhex(SHA256(body))
func HMACAuth(repo domain.APIKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		tsStr := c.GetHeader("X-Timestamp")
		signature := c.GetHeader("X-Signature")

		if apiKey == "" || tsStr == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing required headers: X-API-Key, X-Timestamp, X-Signature",
			})
			return
		}

		// ── Validate timestamp ─────────────────────────────────────────────
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid X-Timestamp"})
			return
		}

		diff := time.Since(time.Unix(ts, 0))
		if diff < 0 {
			diff = -diff
		}
		if diff > timestampTolerance {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "X-Timestamp is too old or too far in the future (±5 min allowed)",
			})
			return
		}

		// ── Read and restore body ──────────────────────────────────────────
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// ── Look up API key ────────────────────────────────────────────────
		key, err := repo.GetByKey(c.Request.Context(), apiKey)
		if err != nil || !key.IsActive() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or revoked API key"})
			return
		}

		// ── Verify signature ───────────────────────────────────────────────
		bodyHash := sha256.Sum256(body)
		message := c.Request.Method + "\n" +
			c.Request.URL.Path + "\n" +
			tsStr + "\n" +
			hex.EncodeToString(bodyHash[:])

		mac := hmac.New(sha256.New, []byte(key.Secret))
		mac.Write([]byte(message))
		expected := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(expected), []byte(signature)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		c.Set("customer_id", key.CustomerID)
		c.Next()
	}
}
