package limiter

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/wychl/limiter/bucket"
	"github.com/wychl/limiter/errors"
	"github.com/wychl/limiter/store"
)

// Response return struct
type Response struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// Limiter config
type Limiter struct {
	rate       float64
	allowDelay bool
	key        func(req *http.Request) string
	Store      store.Store
}

// New return limiter
func New(rate float64) *Limiter {
	return &Limiter{rate: rate, allowDelay: false, key: defaultGetKey}
}

// NewLimiterWithoption  return limiter
func NewLimiterWithoption(rate float64, options ...Option) *Limiter {
	l := &Limiter{rate: rate, allowDelay: false}
	for _, option := range options {
		option(l)
	}
	return l
}

// Handler http middleware
func (l *Limiter) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp Response
		now := time.Now()

		w.Header().Set("Content-Type", "application/json")

		requestKey := l.key(r)
		if requestKey == "" {
			resp.ErrorCode = int(errors.ErrKeyIsNull)
			resp.ErrorMessage = errors.ErrKeyIsNull.String()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		exist, err := l.Store.Exist(requestKey)
		if err != nil {
			resp.ErrorCode = int(errors.ErrStoreRead)
			resp.ErrorMessage = errors.ErrStoreRead.String()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if exist {
			bucket, err := l.Store.Get(requestKey)
			if err != nil {
				resp.ErrorCode = int(errors.ErrStoreRead)
				resp.ErrorMessage = errors.ErrStoreRead.String()
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}

			allow := bucket.Allow(now)
			if !allow {
				resp.ErrorCode = int(errors.ErrToManyRequest)
				resp.ErrorMessage = errors.ErrToManyRequest.String()
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(resp)
				return
			}

			err = l.Store.Set(requestKey, bucket)
			if err != nil {
				resp.ErrorCode = int(errors.ErrStoreSet)
				resp.ErrorMessage = errors.ErrStoreSet.String()
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}

		}

		bucket := &bucket.Bucket{
			AllowDelay: l.allowDelay,
			Timestamp:  now.Unix(),
			Rate:       l.rate,
		}
		err = l.Store.Set(requestKey, bucket)
		if err != nil {
			resp.ErrorCode = int(errors.ErrStoreSet)
			resp.ErrorMessage = errors.ErrStoreSet.String()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		next.ServeHTTP(w, r)
	}
}
