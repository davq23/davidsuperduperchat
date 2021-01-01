package repo

import (
	"context"
	"time"
)

// SessionRepo comprises basic Session operations
type SessionRepo interface {
	// Get gets values from session store
	Get(context.Context, string) (interface{}, error)
	// Set sets values in session store
	Set(context.Context, string, interface{}, time.Duration) (interface{}, error)
	// Delete deletes values in session store
	Delete(context.Context, string) (int64, error)
	// UpdateExpire updates the session expire time
	UpdateExpire(ctx context.Context, sid string, exp time.Duration) error
}
