package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)

	// First 3 requests should be allowed
	for i := 0; i < 3; i++ {
		if !rl.Allow("192.168.1.1") {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// 4th request should be blocked
	if rl.Allow("192.168.1.1") {
		t.Error("4th request should be blocked")
	}

	// Different IP should still be allowed
	if !rl.Allow("192.168.1.2") {
		t.Error("request from different IP should be allowed")
	}
}

func TestRateLimiter_Remaining(t *testing.T) {
	rl := NewRateLimiter(5, time.Minute)

	if remaining := rl.Remaining("192.168.1.1"); remaining != 5 {
		t.Errorf("expected 5 remaining, got %d", remaining)
	}

	rl.Allow("192.168.1.1")
	if remaining := rl.Remaining("192.168.1.1"); remaining != 4 {
		t.Errorf("expected 4 remaining, got %d", remaining)
	}

	rl.Allow("192.168.1.1")
	rl.Allow("192.168.1.1")
	if remaining := rl.Remaining("192.168.1.1"); remaining != 2 {
		t.Errorf("expected 2 remaining, got %d", remaining)
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)

	rl.Allow("192.168.1.1")
	rl.Allow("192.168.1.1")

	if rl.Allow("192.168.1.1") {
		t.Error("should be blocked before reset")
	}

	rl.Reset("192.168.1.1")

	if !rl.Allow("192.168.1.1") {
		t.Error("should be allowed after reset")
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	// Use a very short window for testing
	rl := NewRateLimiter(2, 50*time.Millisecond)

	rl.Allow("192.168.1.1")
	rl.Allow("192.168.1.1")

	if rl.Allow("192.168.1.1") {
		t.Error("should be blocked within window")
	}

	// Wait for window to expire
	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("192.168.1.1") {
		t.Error("should be allowed after window expires")
	}
}

func TestRateLimiter_Concurrent(t *testing.T) {
	rl := NewRateLimiter(100, time.Minute)

	var wg sync.WaitGroup
	allowed := make(chan bool, 200)

	// Launch 200 concurrent requests
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed <- rl.Allow("192.168.1.1")
		}()
	}

	wg.Wait()
	close(allowed)

	// Count allowed requests
	count := 0
	for a := range allowed {
		if a {
			count++
		}
	}

	// Exactly 100 should be allowed
	if count != 100 {
		t.Errorf("expected exactly 100 allowed requests, got %d", count)
	}
}

func TestRateLimiter_MultipleIPs(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)

	// Each IP should have its own limit
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	for _, ip := range ips {
		if !rl.Allow(ip) {
			t.Errorf("first request from %s should be allowed", ip)
		}
		if !rl.Allow(ip) {
			t.Errorf("second request from %s should be allowed", ip)
		}
		if rl.Allow(ip) {
			t.Errorf("third request from %s should be blocked", ip)
		}
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(2, time.Minute)

	r := gin.New()
	r.POST("/login", rl.RateLimit(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// 3rd request should be rate limited
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("3rd request: expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}

func TestRateLimitMiddleware_DifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(1, time.Minute)

	r := gin.New()
	r.POST("/login", rl.RateLimit(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// First IP - should succeed then be blocked
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/login", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("first request from IP1: expected %d, got %d", http.StatusOK, w1.Code)
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/login", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("second request from IP1: expected %d, got %d", http.StatusTooManyRequests, w2.Code)
	}

	// Different IP - should succeed
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/login", nil)
	req3.RemoteAddr = "192.168.1.2:12345"
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("first request from IP2: expected %d, got %d", http.StatusOK, w3.Code)
	}
}

func TestLoginRateLimiter(t *testing.T) {
	rl := LoginRateLimiter(5)

	// Should allow 5 requests
	for i := 0; i < 5; i++ {
		if !rl.Allow("192.168.1.1") {
			t.Errorf("request %d should be allowed (login limiter allows 5)", i+1)
		}
	}

	// 6th should be blocked
	if rl.Allow("192.168.1.1") {
		t.Error("6th request should be blocked by login limiter")
	}
}

func TestLoginRateLimiter_Disabled(t *testing.T) {
	rl := LoginRateLimiter(0)
	if rl != nil {
		t.Error("LoginRateLimiter(0) should return nil")
	}
}

func TestRateLimiter_ZeroRemaining(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)

	rl.Allow("192.168.1.1")
	rl.Allow("192.168.1.1")
	rl.Allow("192.168.1.1") // This one is blocked

	if remaining := rl.Remaining("192.168.1.1"); remaining != 0 {
		t.Errorf("expected 0 remaining, got %d", remaining)
	}
}
