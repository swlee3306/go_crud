package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type RateLimitConfig struct {
	Limit  int           `json:"limit"`  // Number of requests allowed
	Window time.Duration `json:"window"` // Time window for the limit
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	
	// Start cleanup goroutine
	go rl.cleanup()
	
	return rl
}

func (rl *RateLimiter) IsAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	// Get existing requests for this key
	requests, exists := rl.requests[key]
	if !exists {
		requests = []time.Time{}
	}
	
	// Remove old requests outside the window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if we're under the limit
	if len(validRequests) < rl.limit {
		// Add current request
		validRequests = append(validRequests, now)
		rl.requests[key] = validRequests
		return true
	}
	
	// Update requests (without adding current request)
	rl.requests[key] = validRequests
	return false
}

func (rl *RateLimiter) GetRemainingRequests(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	requests, exists := rl.requests[key]
	if !exists {
		return rl.limit
	}
	
	// Count valid requests
	validCount := 0
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validCount++
		}
	}
	
	remaining := rl.limit - validCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (rl *RateLimiter) GetResetTime(key string) time.Time {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	
	requests, exists := rl.requests[key]
	if !exists || len(requests) == 0 {
		return time.Now()
	}
	
	// Find the oldest request
	oldest := requests[0]
	for _, reqTime := range requests {
		if reqTime.Before(oldest) {
			oldest = reqTime
		}
	}
	
	// Reset time is when the oldest request will expire
	return oldest.Add(rl.window)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)
		
		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			
			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

// Rate limiting middleware
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address or user ID)
		clientID := getClientID(c)
		
		// Check if request is allowed
		if !rl.IsAllowed(clientID) {
			remaining := rl.GetRemainingRequests(clientID)
			resetTime := rl.GetResetTime(clientID)
			
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			c.Abort()
			return
		}
		
		// Add rate limit headers
		remaining := rl.GetRemainingRequests(clientID)
		resetTime := rl.GetResetTime(clientID)
		
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
		
		c.Next()
	}
}

// IP-based rate limiting
func IPRateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !rl.IsAllowed(clientIP) {
			remaining := rl.GetRemainingRequests(clientIP)
			resetTime := rl.GetResetTime(clientIP)
			
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests from this IP, please try again later",
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			c.Abort()
			return
		}
		
		remaining := rl.GetRemainingRequests(clientIP)
		resetTime := rl.GetResetTime(clientIP)
		
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
		
		c.Next()
	}
}

// User-based rate limiting
func UserRateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// If no user ID, fall back to IP-based limiting
			clientIP := c.ClientIP()
			if !rl.IsAllowed(clientIP) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}
		
		userKey := fmt.Sprintf("user:%v", userID)
		
		if !rl.IsAllowed(userKey) {
			remaining := rl.GetRemainingRequests(userKey)
			resetTime := rl.GetResetTime(userKey)
			
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			c.Abort()
			return
		}
		
		remaining := rl.GetRemainingRequests(userKey)
		resetTime := rl.GetResetTime(userKey)
		
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
		
		c.Next()
	}
}

func getClientID(c *gin.Context) string {
	// Try to get user ID first
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%v", userID)
	}
	
	// Fall back to IP address
	return c.ClientIP()
}

// Global rate limiter instances
var (
	// General API rate limiter: 100 requests per minute
	GeneralRateLimiter = NewRateLimiter(100, time.Minute)
	
	// Auth rate limiter: 5 requests per minute
	AuthRateLimiter = NewRateLimiter(5, time.Minute)
	
	// Strict rate limiter: 10 requests per minute
	StrictRateLimiter = NewRateLimiter(10, time.Minute)
)
