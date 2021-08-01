package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit limits rate.
func RateLimit() func(*gin.Context) {
	var mu sync.Mutex
	m := map[string]*rate.Limiter{}

	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		var limiter *rate.Limiter
		var ok bool
		ip := c.ClientIP()

		if limiter, ok = m[ip]; !ok {
			limiter = rate.NewLimiter(rate.Every(100*time.Millisecond), 20)
		}

		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
		}

		m[ip] = limiter
	}
}
