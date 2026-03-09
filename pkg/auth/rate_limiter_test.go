package auth

import (
	"fmt"
	"sync"
	"testing"
)

func TestRateLimiter_AllowsWithinBurst(t *testing.T) {
	rl := NewRateLimiter(10, 5)
	key := "bk_testkey"

	for i := 0; i < 5; i++ {
		if !rl.get(key).Allow() {
			t.Errorf("request %d should be allowed within burst of 5", i+1)
		}
	}
}

func TestRateLimiter_BlocksAfterBurstExhausted(t *testing.T) {
	// Near-zero refill rate so tokens won't replenish during the test.
	rl := NewRateLimiter(0.0001, 3)
	key := "bk_testkey"

	for i := 0; i < 3; i++ {
		rl.get(key).Allow()
	}

	if rl.get(key).Allow() {
		t.Error("expected request to be blocked after burst exhausted")
	}
}

func TestRateLimiter_PerKeyIsolation(t *testing.T) {
	rl := NewRateLimiter(0.0001, 2)

	// Exhaust key A's bucket.
	rl.get("bk_keyA").Allow()
	rl.get("bk_keyA").Allow()

	// Key B must have its own independent bucket.
	if !rl.get("bk_keyB").Allow() {
		t.Error("key B should be unaffected by key A's rate limit")
	}
}

func TestRateLimiter_NewKeyGetsFullBurst(t *testing.T) {
	rl := NewRateLimiter(10, 10)

	// A brand-new key should start with a full burst.
	limiter := rl.get("bk_newkey")
	if !limiter.Allow() {
		t.Error("first request for a new key should always be allowed")
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(1000, 1000)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("bk_key_%d", i%10)
			rl.get(key).Allow()
		}(i)
	}

	wg.Wait() // must not race or deadlock
}

func TestRateLimiter_SameKeyReturnsSameLimiter(t *testing.T) {
	rl := NewRateLimiter(10, 10)

	l1 := rl.get("bk_same")
	l2 := rl.get("bk_same")

	if l1 != l2 {
		t.Error("get() must return the same limiter for the same key")
	}
}
