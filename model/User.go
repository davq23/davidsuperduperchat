package model

// User represents an user in the db
type User struct {
	ID       int
	Username string
	Hash     string
}
