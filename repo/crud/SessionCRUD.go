package crud

import (
	"context"
	"davidws/utils/channels"
	"time"

	"github.com/go-redis/redis/v8"
)

// SessionCRUD allows basic Session operations
type SessionCRUD struct {
	redis *redis.Client
}

// NewSessionCRUD creates a *SessionCRUD
func NewSessionCRUD(redis *redis.Client) *SessionCRUD {
	return &SessionCRUD{
		redis: redis,
	}
}

// Get gets values from session store
func (sr *SessionCRUD) Get(ctx context.Context, sid string) (val interface{}, err error) {

	done := make(chan bool)

	go func(done chan<- bool) {
		val, err = sr.redis.Get(ctx, sid).Result()
		done <- true
	}(done)

	if channels.OK(done) {

	}

	return
}

// UpdateExpire value
func (sr *SessionCRUD) UpdateExpire(ctx context.Context, sid string, exp time.Duration) (err error) {
	done := make(chan bool)

	go func(done chan<- bool) {
		_, err = sr.redis.Expire(ctx, sid, exp).Result()
		done <- true
	}(done)

	if channels.OK(done) {

	}

	return
}

// Set sets values in session store
func (sr *SessionCRUD) Set(ctx context.Context, sid string, data interface{}, exp time.Duration) (val interface{}, err error) {
	done := make(chan bool)

	go func(done chan<- bool) {
		val, err = sr.redis.Set(ctx, sid, data, exp).Result()
		done <- true
	}(done)

	if channels.OK(done) {

	}

	return
}

// Delete deletes values in session store
func (sr *SessionCRUD) Delete(ctx context.Context, sid string) (val int64, err error) {
	done := make(chan bool)

	go func(done chan<- bool) {
		val, err = sr.redis.Del(ctx, sid).Result()
		done <- true
	}(done)

	if channels.OK(done) {

	}

	return
}
