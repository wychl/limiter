package store

import (
	"sync"

	"github.com/wychl/limiter/bucket"
	"github.com/wychl/limiter/errors"

	"github.com/garyburd/redigo/redis"
)

// Redis redis store bucket
type Redis struct {
	mutex *sync.Mutex
	pool  *redis.Pool
}

var _ Store = &Redis{}

// NewRedis create redis store
func NewRedis(pool *redis.Pool) *Redis {
	conn := pool.Get()
	defer conn.Close()

	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		panic(err)
	}

	return &Redis{pool: pool, mutex: new(sync.Mutex)}
}

// Set store bucket to redis
func (r *Redis) Set(key string, b *bucket.Bucket) error {
	if key == "" {
		return errors.ErrKeyIsNull
	}

	if b == nil {
		return errors.ErrInvalidBucket
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	conn := r.pool.Get()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send("HSET", key, allowDelayField, b.AllowDelay)
	if err != nil {
		return err
	}

	err = conn.Send("HSET", key, rateField, b.Rate)
	if err != nil {
		return err
	}

	err = conn.Send("HSET", key, timestampField, b.Timestamp)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return conn.Close()
}

// Get get bucket from redis
func (r *Redis) Get(key string) (*bucket.Bucket, error) {
	if key == "" {
		return nil, errors.ErrKeyIsNull
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	conn := r.pool.Get()
	defer conn.Close()

	var (
		b   bucket.Bucket
		err error
	)

	b.AllowDelay, err = redis.Bool(conn.Do("HGET", key, allowDelayField))
	if err != nil {
		return nil, err
	}

	b.Rate, err = redis.Float64(conn.Do("HGET", key, rateField))
	if err != nil {
		return nil, err
	}

	b.Timestamp, err = redis.Int64(conn.Do("HGET", key, timestampField))
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// Exist key exit in redis
func (r *Redis) Exist(key string) (bool, error) {
	if key == "" {
		return false, errors.ErrKeyIsNull
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	conn := r.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXIST", key))
}
