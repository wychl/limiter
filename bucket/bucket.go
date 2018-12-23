package bucket

import (
	"time"
)

// Bucket bucket config
type Bucket struct {
	Timestamp  int64
	Rate       float64
	AllowDelay bool
}

// Allow allow current request
func (b *Bucket) Allow(now time.Time) bool {
	internal := float64(now.Unix()-b.Timestamp) * b.Rate
	if internal < 1 {
		return false
	}

	if !b.AllowDelay {
		b.Timestamp = now.Unix()
	} else {
		b.Timestamp = b.Timestamp + int64(1/b.Rate)
	}

	return true
}
