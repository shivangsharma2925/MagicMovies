package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	rate    rate.Limit
	burst   int
}

// constructor
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rate:    r,
		burst:   b,
	}

	go rl.cleanup()

	return rl
}

// cleanup old clients, the cleanup will be done once every 1 minute where map is locked since map is not thread-safe.
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for key, c := range rl.clients {
			if time.Since(c.lastSeen) > 5*time.Minute {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

// get limiter for a key (IP or userID)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if c, exists := rl.clients[key]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)

	rl.clients[key] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

func (rl *RateLimiter) IPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		ip := c.ClientIP()

		limiter := rl.getLimiter(ip)

		if !limiter.Allow(){
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		userId, exists := c.Get("userId")
		if !exists{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		id, ok := userId.(string)

		if !ok{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User"})
			c.Abort()
			return
		}

		limiter := rl.getLimiter(id)

		if !limiter.Allow(){
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
