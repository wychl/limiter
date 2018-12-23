package store

import (
	"sync"

	"github.com/wychl/limiter/bucket"
	"github.com/wychl/limiter/errors"
)

// Memory memory store bucket
type Memory struct {
	mutex  *sync.Mutex
	memory map[string]*bucket.Bucket
}

var _ Store = &Memory{}

// NewMemory create memory store
func NewMemory(memory map[string]*bucket.Bucket) *Memory {
	return &Memory{memory: memory, mutex: new(sync.Mutex)}
}

// Set store bucket to memory
func (r *Memory) Set(key string, b *bucket.Bucket) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.memory[key] = b

	return nil
}

// Get get bucket from memory
func (r *Memory) Get(key string) (*bucket.Bucket, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	b, exist := r.memory[key]
	if !exist {
		return nil, errors.ErrNoExist
	}

	return b, nil
}

// Exist key exit in memory
func (r *Memory) Exist(key string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exist := r.memory[key]
	return exist, nil
}
