package auth_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/pkg/auth"
	"github.com/belayhun-arage/billing-service/test/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// sign builds the HMAC-SHA256 signature using the same algorithm as the middleware.
func sign(method, path, tsStr string, body []byte, secret string) string {
	bodyHash := sha256.Sum256(body)
	message := method + "\n" + path + "\n" + tsStr + "\n" + hex.EncodeToString(bodyHash[:])
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// newRouter returns a minimal Gin router with HMACAuth applied.
func newRouter(repo domain.APIKeyRepository) *gin.Engine {
	r := gin.New()
	r.Use(auth.HMACAuth(repo))
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"api_key_id":    c.GetString("api_key_id"),
			"api_key_label": c.GetString("api_key_label"),
		})
	})
	return r
}

// fixture builds a consistent test API key and a matching mock repo.
func fixture() (*domain.APIKey, domain.APIKeyRepository) {
	key := &domain.APIKey{
		ID:     "key-id-1",
		Key:    "bk_testfixture",
		Secret: "supersecrethmackey1234567890abcdef",
		Label:  "test",
	}
	repo := &mocks.MockAPIKeyRepository{
		GetByKeyFn: func(_ context.Context, k string) (*domain.APIKey, error) {
			if k == key.Key {
				return key, nil
			}
			return nil, domain.ErrAPIKeyNotFound
		},
	}
	return key, repo
}

func TestHMACAuth_ValidRequest(t *testing.T) {
	key, repo := fixture()
	body := []byte(`{"amount":5000}`)
	ts := fmt.Sprintf("%d", time.Now().Unix())
	sig := sign("POST", "/test", ts, body, key.Secret)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("X-API-Key", key.Key)
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", sig)

	newRouter(repo).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify api_key_id and api_key_label are propagated into the handler context.
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["api_key_id"] != key.ID {
		t.Errorf("api_key_id = %q, want %q", resp["api_key_id"], key.ID)
	}
	if resp["api_key_label"] != key.Label {
		t.Errorf("api_key_label = %q, want %q", resp["api_key_label"], key.Label)
	}
}

func TestHMACAuth_MissingHeaders(t *testing.T) {
	_, repo := fixture()
	r := newRouter(repo)

	tests := []struct {
		desc    string
		headers map[string]string
	}{
		{
			"missing all headers",
			map[string]string{},
		},
		{
			"missing X-Timestamp and X-Signature",
			map[string]string{"X-API-Key": "bk_testfixture"},
		},
		{
			"missing X-Signature",
			map[string]string{
				"X-API-Key":   "bk_testfixture",
				"X-Timestamp": fmt.Sprintf("%d", time.Now().Unix()),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{}`)))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			r.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", w.Code)
			}
		})
	}
}

func TestHMACAuth_InvalidTimestamp(t *testing.T) {
	key, repo := fixture()
	body := []byte(`{}`)

	tests := []struct {
		desc string
		ts   string
	}{
		{"not a number", "notanumber"},
		{"stale timestamp (too old)", fmt.Sprintf("%d", time.Now().Add(-2*time.Minute).Unix())},
		{"future timestamp (too far ahead)", fmt.Sprintf("%d", time.Now().Add(2*time.Minute).Unix())},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			sig := sign("POST", "/test", tc.ts, body, key.Secret)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
			req.Header.Set("X-API-Key", key.Key)
			req.Header.Set("X-Timestamp", tc.ts)
			req.Header.Set("X-Signature", sig)

			newRouter(repo).ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("%s: expected 401, got %d", tc.desc, w.Code)
			}
		})
	}
}

func TestHMACAuth_InvalidAPIKey(t *testing.T) {
	_, repo := fixture()
	body := []byte(`{}`)
	ts := fmt.Sprintf("%d", time.Now().Unix())
	sig := sign("POST", "/test", ts, body, "anysecret")

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("X-API-Key", "bk_doesnotexist")
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", sig)

	newRouter(repo).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unknown API key, got %d", w.Code)
	}
}

func TestHMACAuth_RevokedKey(t *testing.T) {
	now := time.Now()
	revokedKey := &domain.APIKey{
		ID:        "key-revoked",
		Key:       "bk_revoked",
		Secret:    "revokedsecret1234567890abcdef1234",
		Label:     "revoked-key",
		RevokedAt: &now,
	}
	repo := &mocks.MockAPIKeyRepository{
		GetByKeyFn: func(_ context.Context, k string) (*domain.APIKey, error) {
			return revokedKey, nil
		},
	}

	body := []byte(`{}`)
	ts := fmt.Sprintf("%d", time.Now().Unix())
	sig := sign("POST", "/test", ts, body, revokedKey.Secret)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("X-API-Key", revokedKey.Key)
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", sig)

	newRouter(repo).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for revoked key, got %d", w.Code)
	}
}

func TestHMACAuth_WrongSignature(t *testing.T) {
	key, repo := fixture()
	body := []byte(`{"amount":5000}`)
	ts := fmt.Sprintf("%d", time.Now().Unix())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("X-API-Key", key.Key)
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", "deadbeefdeadbeefdeadbeefdeadbeef")

	newRouter(repo).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong signature, got %d", w.Code)
	}
}

func TestHMACAuth_SignatureIsMethodSensitive(t *testing.T) {
	key, repo := fixture()
	body := []byte(`{}`)
	ts := fmt.Sprintf("%d", time.Now().Unix())

	// Sign as GET but send as POST — the middleware must reject it.
	sig := sign("GET", "/test", ts, body, key.Secret)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("X-API-Key", key.Key)
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", sig)

	newRouter(repo).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when method in signature doesn't match, got %d", w.Code)
	}
}

func TestHMACAuth_SignatureIsBodySensitive(t *testing.T) {
	key, repo := fixture()
	ts := fmt.Sprintf("%d", time.Now().Unix())

	// Sign an empty body but send a non-empty one.
	sig := sign("POST", "/test", ts, []byte(`{}`), key.Secret)
	tamperedBody := []byte(`{"amount":9999}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(tamperedBody))
	req.Header.Set("X-API-Key", key.Key)
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", sig)

	newRouter(repo).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when body tampered after signing, got %d", w.Code)
	}
}
