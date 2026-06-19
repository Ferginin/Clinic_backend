package middleware

import (
	"net/http"
	"sync"
	"time"

	"Clinic_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type windowEntry struct {
	mu          sync.Mutex
	count       int
	windowStart time.Time
}

type ipRateLimiter struct {
	mu      sync.Mutex
	entries map[string]*windowEntry
	limit   int
	window  time.Duration
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	rl := &ipRateLimiter{
		entries: make(map[string]*windowEntry),
		limit:   limit,
		window:  window,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.window * 2)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		rl.mu.Lock()
		for ip, e := range rl.entries {
			e.mu.Lock()
			if now.Sub(e.windowStart) > rl.window*2 {
				delete(rl.entries, ip)
			}
			e.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func (rl *ipRateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	e, ok := rl.entries[ip]
	if !ok {
		e = &windowEntry{windowStart: time.Now()}
		rl.entries[ip] = e
	}
	rl.mu.Unlock()

	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	if now.Sub(e.windowStart) > rl.window {
		e.count = 0
		e.windowStart = now
	}
	e.count++
	return e.count <= rl.limit
}

// RateLimiterMiddleware limits requests per IP within a sliding window.
// limit — max requests, window — time window duration.
func RateLimiterMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rl := newIPRateLimiter(limit, window)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "too many requests, please try again later")
			c.Abort()
			return
		}
		c.Next()
	}
}
