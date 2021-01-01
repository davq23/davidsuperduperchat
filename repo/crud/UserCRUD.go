package crud

import (
	"context"

	"davidws/model"
	"davidws/utils/channels"

	"github.com/jackc/pgx/v4/pgxpool"
)

// UserCRUD allows user db operations
type UserCRUD struct {
	userPool *pgxpool.Pool
}

// NewUserCRUD creates a *UserCRUD
func NewUserCRUD(userPool *pgxpool.Pool) *UserCRUD {
	return &UserCRUD{
		userPool: userPool,
	}
}

// GetByID returns an User according to its ID
func (uc *UserCRUD) GetByID(ctx context.Context, userID string) (user model.User, err error) {
	done := make(chan bool)

	go func(c chan<- bool) {
		row := uc.userPool.QueryRow(ctx, "SELECT username, hash FRoM users WHERE user_id = $1", userID)

		err = row.Scan(&user.Username, &user.Hash)

		if err != nil {
			c <- false
			return
		}

		c <- true
	}(done)

	if channels.OK(done) {
		user.ID = userID
	}

	return
}

// GetByUsername returns an User according to its username
func (uc *UserCRUD) GetByUsername(ctx context.Context, username string) (user model.User, err error) {
	done := make(chan bool)

	go func(c chan<- bool) {
		row := uc.userPool.QueryRow(ctx, "SELECT user_id, hash FRoM users WHERE username = $1", username)

		err = row.Scan(&user.ID, &user.Hash)

		if err != nil {
			c <- false
			return
		}

		c <- true
	}(done)

	if channels.OK(done) {
		user.Username = username
	}

	return
}

// Insert inserts one user
func (uc *UserCRUD) Insert(ctx context.Context, user model.User) (err error) {
	done := make(chan bool)

	go func(c chan<- bool) {
		row := uc.userPool.QueryRow(ctx, "INSERT INTo users (username, hash) VALUES ($1, $2) RETURNING user_id", user.Username, user.Hash)

		err = row.Scan(&user.ID)

		if err != nil {
			c <- false
			return
		}

		c <- true
	}(done)

	if channels.OK(done) {
	}

	return
}
