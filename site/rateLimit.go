package site

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	siteName    string
	minuteLimit *TokenBucket
	hourLimit   *TokenBucket
	dayLimit    *TokenBucket
	lock        sync.Mutex
}

type TokenBucket struct {
	tokens   float64
	capacity float64
	rate     float64
	lastTime time.Time
	timeUnit time.Duration
}

func NewTokenBucket(capacity float64, timeUnit time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:   capacity,
		capacity: capacity,
		rate:     capacity / float64(timeUnit),
		lastTime: time.Now(),
		timeUnit: timeUnit,
	}
}

func (tb *TokenBucket) Allow() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastTime).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)
	tb.lastTime = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

func NewRateLimiter(siteName string, maxPerMin, maxPerHour, maxPerDay int) *RateLimiter {
	return &RateLimiter{
		siteName:    siteName,
		minuteLimit: NewTokenBucket(float64(maxPerMin), time.Minute),
		hourLimit:   NewTokenBucket(float64(maxPerHour), time.Hour),
		dayLimit:    NewTokenBucket(float64(maxPerDay), 24*time.Hour),
		lock:        sync.Mutex{},
	}
}

func (rl *RateLimiter) Allow() error {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if !rl.minuteLimit.Allow() {
		return fmt.Errorf("站点 %s 每分钟请求次数超限", rl.siteName)
	}
	if !rl.hourLimit.Allow() {
		return fmt.Errorf("站点 %s 每小时请求次数超限", rl.siteName)
	}
	if !rl.dayLimit.Allow() {
		return fmt.Errorf("站点 %s 每天请求次数超限", rl.siteName)
	}
	return nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
