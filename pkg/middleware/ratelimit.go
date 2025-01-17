package middleware

import (
	"example-go-project/pkg/utils"
	"net/http"
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
		// user, ok := GetUserFromContext(c)
		// if !ok {
		// 	utils.SendError(c, http.StatusUnauthorized, "User not found")
		// 	return
		// }

		if !limiter.Allow(key) {
			utils.SendError(c, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
