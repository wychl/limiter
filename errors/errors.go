package errors

// LimiterErr  limiter error type
type LimiterErr int

const (
	// ErrKeyIsNull key is null
	ErrKeyIsNull LimiterErr = iota + 1000
	// ErrStoreRead read error
	ErrStoreRead
	// ErrStoreSet set error
	ErrStoreSet
	// ErrToManyRequest too many request
	ErrToManyRequest
	// ErrNoExist key not in store
	ErrNoExist
	//ErrInvalidBucket bucket invalid
	ErrInvalidBucket
)

var limiterErrorMessage = map[LimiterErr]string{
	ErrKeyIsNull:     "key is null",
	ErrStoreRead:     "store read error",
	ErrToManyRequest: "too many request",
	ErrStoreSet:      "store set error",
	ErrNoExist:       "key not exist",
	ErrInvalidBucket: "bucket is invalid",
}

// String return string
func (l LimiterErr) String() string {
	message, ok := limiterErrorMessage[l]
	if !ok {
		return "unknown error"
	}
	return message
}

// Error  error interface
func (l LimiterErr) Error() string {
	return l.String()
}
