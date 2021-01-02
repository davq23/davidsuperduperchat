package repo

import (
	"context"
	"davidws/model"
)

// UserRepo  comprises basic User operations
type UserRepo interface {
	// GetByID returns an User according to its ID
	GetByID(ctx context.Context, userID int) (model.User, error)
	// GetByUsername returns an User according to its username
	GetByUsername(ctx context.Context, username string) (user model.User, err error)
	// Insert inserts a new User and returns an error if unsuccessful
	Insert(ctx context.Context, user model.User) error
}
