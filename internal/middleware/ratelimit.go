package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	rate     int
	interval time.Duration
	mu       sync.Mutex
	tokens   map[string][]time.Time
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		interval: interval,
		tokens:   make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.interval)

	if _, exists := rl.tokens[key]; !exists {
		rl.tokens[key] = []time.Time{now}
		return true
	}

	var validTokens []time.Time
	for _, t := range rl.tokens[key] {
		if t.After(windowStart) {
			validTokens = append(validTokens, t)
		}
	}

	if len(validTokens) < rl.rate {
		validTokens = append(validTokens, now)
		rl.tokens[key] = validTokens
		return true
	}

	rl.tokens[key] = validTokens
	return false
}

func RateLimit(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	return func(c *gin.Context) {
		// use ip
		key := c.ClientIP()

		// or use user id on jwt
		// if user, exists := c.Get("user"); exists {
		//     key = user.(string) // แปลง user ID เป็น string
		// }

		if !limiter.Allow(key) {
			c.JSON(429, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}