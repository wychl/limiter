package store

import "github.com/wychl/limiter/bucket"

// Store store interface
type Store interface {
	Set(key string, bucket *bucket.Bucket) error
	Get(key string) (*bucket.Bucket, error)
	Exist(key string) (bool, error)
}
