package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/AFR-projection/GTA/backend/internal/httpx"
)

type visitor struct {
	count    int
	windowStart time.Time
}

// RateLimit is a simple per-IP fixed window limiter (no Redis required).
func RateLimit(maxPerMinute int) func(http.Handler) http.Handler {
	if maxPerMinute <= 0 {
		maxPerMinute = 120
	}
	var mu sync.Mutex
	visitors := map[string]*visitor{}

	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			mu.Lock()
			now := time.Now()
			for ip, v := range visitors {
				if now.Sub(v.windowStart) > 2*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			now := time.Now()

			mu.Lock()
			v, ok := visitors[ip]
			if !ok || now.Sub(v.windowStart) >= time.Minute {
				visitors[ip] = &visitor{count: 1, windowStart: now}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}
			v.count++
			if v.count > maxPerMinute {
				mu.Unlock()
				w.Header().Set("Retry-After", "60")
				httpx.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
