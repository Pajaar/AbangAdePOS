package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit membatasi request per-IP (token bucket), dipakai khusus di
// route login. Cleanup berkala membuang IP yang sudah lama tidak nongol
// supaya map tidak bocor memori (cukup untuk single instance).
func RateLimit(perMinute float64, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)
	ratePerSecond := rate.Limit(perMinute / 60)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, ok := visitors[ip]
		if !ok {
			v = &visitor{limiter: rate.NewLimiter(ratePerSecond, burst)}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		allowed := v.limiter.Allow()
		mu.Unlock()

		if !allowed {
			util.Error(c, apperror.TooManyRequests("Terlalu banyak percobaan, coba lagi nanti"))
			c.Abort()
			return
		}

		c.Next()
	}
}
