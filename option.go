package limiter

import (
	"net/http"

	"github.com/wychl/limiter/store"

	"github.com/rs/xid"
)

// Option limiter middleware
type Option func(l *Limiter) *Limiter

// Key access key from request
type Key func(req *http.Request) string

func defaultGetKey(req *http.Request) string {
	clientIP := getClientIP(req)
	if clientIP != "" {
		return "limiter::" + clientIP
	}

	return "limiter::" + xid.New().String()
}

// KeyOption return limiter key middleware
func KeyOption(key Key) Option {
	return func(l *Limiter) *Limiter {
		l.key = key
		return l
	}
}

// AllowDelayOption return limiter delay middleware
func AllowDelayOption(allowDelay bool) Option {
	return func(l *Limiter) *Limiter {
		l.allowDelay = allowDelay
		return l
	}
}

// StoreOption define limiter store
func StoreOption(s store.Store) Option {
	return func(l *Limiter) *Limiter {
		l.store = s
		return l
	}
}
