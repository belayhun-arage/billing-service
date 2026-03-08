package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter holds a per-key token bucket store.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	r       rate.Limit
	burst   int
}

// NewRateLimiter creates a RateLimiter with rps tokens per second and the
// given burst size. A background goroutine removes stale entries every 5 min.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*entry),
		r:       rate.Limit(rps),
		burst:   burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) get(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	e, ok := rl.entries[key]
	if !ok {
		e = &entry{limiter: rate.NewLimiter(rl.r, rl.burst)}
		rl.entries[key] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

// cleanup removes entries that haven't been seen in 10 minutes.
func (rl *RateLimiter) cleanup() {
	for range time.Tick(5 * time.Minute) {
		rl.mu.Lock()
		for k, e := range rl.entries {
			if time.Since(e.lastSeen) > 10*time.Minute {
				delete(rl.entries, k)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit returns a Gin middleware that enforces per-API-key rate limiting.
// Must run after HMACAuth so X-API-Key is guaranteed to be present and valid.
func RateLimit(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		limiter := rl.get(key)

		if !limiter.Allow() {
			retryAfter := int(time.Second/time.Duration(rl.r)) + 1
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("rate limit exceeded, retry after %ds", retryAfter),
			})
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
		c.Next()
	}
}
